package signer

import (
	"fmt"
	"github.com/lulucas/go-signer/hash"
	"reflect"
	"sort"
	"strings"
)

type (
	KvJoinFn   func(key, value string) string
	PostHookFn func(s, joinChar, key string, kvJoinFn KvJoinFn) string
)

// keep key even if it has a empty value
func (s *Signer) NoSkipEmpty() *Signer {
	s.skipEmpty = false
	return s
}

//  Func to sign str
func (s *Signer) HashFunc(fn hash.Func) *Signer {
	s.hashFunc = fn
	return s
}

// keys that do not be signed
func (s *Signer) IgnoreKeys(keys ...string) *Signer {
	s.ignoreKeys = map[string]struct{}{}
	for _, key := range keys {
		s.ignoreKeys[key] = struct{}{}
	}
	return s
}

// key name from struct tag, when input type is struct
func (s *Signer) Tag(tag string) *Signer {
	s.tag = tag
	return s
}

// key add to str
func (s *Signer) Key(key string) *Signer {
	s.key = key
	return s
}

// join char
func (s *Signer) JoinChar(char string) *Signer {
	s.joinChar = char
	return s
}

// function that join key and value
func (s *Signer) KvJoinFunc(fn KvJoinFn) *Signer {
	s.kvJoinFunc = fn
	return s
}

func (s *Signer) PostHookFunc(fn PostHookFn) *Signer {
	s.postHookFunc = fn
	return s
}

type Signer struct {
	key          string
	ignoreKeys   map[string]struct{}
	skipEmpty    bool
	joinChar     string
	kvJoinFunc   KvJoinFn
	hashFunc     hash.Func
	postHookFunc PostHookFn
	tag          string
}

func New() *Signer {
	s := &Signer{
		ignoreKeys: map[string]struct{}{"sign": {}},
		skipEmpty:  true,
		joinChar:   "&",
		kvJoinFunc: func(key, value string) string {
			return key + "=" + value
		},
		hashFunc: hash.MD5(false),
	}
	s.postHookFunc = func(s, joinChar, key string, kvJoinFunc KvJoinFn) string {
		return s + joinChar + kvJoinFunc("key", key)
	}
	return s
}

// Sign str and get result
func (s *Signer) Sign(data interface{}) string {
	str := s.StrToSign(data)
	return s.hashFunc(str)
}

// Get str to sign, can be used for debugging
func (s *Signer) StrToSign(data interface{}) string {
	m := map[string]string{}
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			key := t.Field(i).Name
			if s.tag != "" {
				key = strings.Split(t.Field(i).Tag.Get(s.tag), ",")[0]
			}
			m[key] = fmt.Sprintf("%v", val.Field(i).Interface())
		}
	case reflect.Map:
		for _, element := range val.MapKeys() {
			if val.MapIndex(element).Kind() == reflect.Slice {
				if val.MapIndex(element).Len() > 0 {
					m[element.String()] = fmt.Sprintf("%v", val.MapIndex(element).Index(0).Interface())
				}
			} else {
				m[element.String()] = fmt.Sprintf("%v", val.MapIndex(element).Interface())
			}
		}
	case reflect.Ptr:
		return s.StrToSign(reflect.ValueOf(data).Elem().Interface())
	}

	var keys []string
	for k, v := range m {
		_, skip := s.ignoreKeys[k]
		if (s.skipEmpty && v == "") || skip {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairsToSign []string
	for _, k := range keys {
		pairsToSign = append(pairsToSign, s.kvJoinFunc(k, m[k]))
	}

	strToSign := strings.Join(pairsToSign, s.joinChar)

	return s.postHookFunc(strToSign, s.joinChar, s.key, s.kvJoinFunc)
}
