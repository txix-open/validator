package ar

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/txix-open/validator/v10"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {
	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "required",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_if",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_unless",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_with",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_with_all",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_without",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "required_without_all",
			translation: "حقل {0} مطلوب",
			override:    false,
		},
		{
			tag:         "excluded_if",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "excluded_unless",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "excluded_with",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "excluded_with_all",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "excluded_without",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "excluded_without_all",
			translation: "حقل {0} مستبعد",
			override:    false,
		},
		{
			tag:         "isdefault",
			translation: "حقل {0} يجب أن يكون قيمة إفتراضية",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "يجب أن يكون طول {0} مساويا ل {1}", false); err != nil {
					return
				}
				if err = ut.Add("len-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("len-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("len-number", "يجب أن يكون {0} مساويا ل {1}", false); err != nil {
					return
				}
				if err = ut.Add("len-items", "يجب أن يحتوي {0} على {1}", false); err != nil {
					return
				}
				if err = ut.Add("len-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("len-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("len-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("len-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("len-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("len-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("len-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("len-items", fe.Field(), c)

				default:
					t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("min-string", "{0} يجب أن يكون {1} على الأقل", false); err != nil {
					return
				}
				if err = ut.Add("min-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("min-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("min-number", "{0} يجب أن يكون {1} أو أكثر", false); err != nil {
					return
				}
				if err = ut.Add("min-items", "يجب أن يحتوي {0} على {1} على الأقل", false); err != nil {
					return
				}
				if err = ut.Add("min-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("min-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("min-duration", "يجب أن تكون مدة {0} {1} أو أكبر", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if fe.Type() != reflect.TypeOf(time.Duration(0)) {
						return nil
					}

					t, err = ut.T("min-duration", fe.Field(), fe.Param())
					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				if err = fn(); err != nil || t != "" {
					goto END
				}

				if fe.Type() == reflect.TypeOf(time.Duration(0)) {
					t, err = ut.T("min-duration", fe.Field(), fe.Param())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("min-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("min-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("min-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("min-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("min-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("min-items", fe.Field(), c)

				default:
					t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("max-string", "يجب أن يكون طول {0} بحد أقصى {1}", false); err != nil {
					return
				}
				if err = ut.Add("max-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("max-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("max-number", "{0} يجب أن يكون {1} أو اقل", false); err != nil {
					return
				}
				if err = ut.Add("max-items", "يجب أن يحتوي {0} على {1} كحد أقصى", false); err != nil {
					return
				}
				if err = ut.Add("max-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("max-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("max-duration", "يجب أن تكون مدة {0} {1} أو أقل", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if fe.Type() != reflect.TypeOf(time.Duration(0)) {
						return nil
					}

					t, err = ut.T("max-duration", fe.Field(), fe.Param())
					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				if err = fn(); err != nil || t != "" {
					goto END
				}

				if fe.Type() == reflect.TypeOf(time.Duration(0)) {
					t, err = ut.T("max-duration", fe.Field(), fe.Param())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("max-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("max-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("max-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("max-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("max-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("max-items", fe.Field(), c)

				default:
					t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag:         "eq",
			translation: "{0} لا يساوي {1}",
			override:    false,
		},
		{
			tag:         "ne",
			translation: "{0} يجب ألا يساوي {1}",
			override:    false,
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("lt-string", "يجب أن يكون طول {0} أقل من {1}", false); err != nil {
					return
				}
				if err = ut.Add("lt-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("lt-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("lt-number", "يجب أن يكون {0} أقل من {1}", false); err != nil {
					return
				}
				if err = ut.Add("lt-items", "يجب أن يحتوي {0} على أقل من {1}", false); err != nil {
					return
				}
				if err = ut.Add("lt-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("lt-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("lt-datetime", "يجب أن يكون {0} أقل من التاريخ والوقت الحاليين", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var digits uint64
				var kind reflect.Kind
				var f64 float64

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				// Special handling for time.Time
				if fe.Type() == reflect.TypeOf(time.Time{}) {
					t, err = ut.T("lt-datetime", fe.Field())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("lt-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("lt-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("lt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("lt-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("lt-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("lt-items", fe.Field(), c)

				case reflect.Struct:
					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))

				default:
					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("lte-string", "يجب أن يكون طول {0} كحد أقصى {1}", false); err != nil {
					return
				}
				if err = ut.Add("lte-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("lte-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("lte-number", "{0} يجب أن يكون {1} أو اقل", false); err != nil {
					return
				}
				if err = ut.Add("lte-items", "يجب أن يحتوي {0} على {1} كحد أقصى", false); err != nil {
					return
				}
				if err = ut.Add("lte-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("lte-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("lte-datetime", "يجب أن يكون {0} أقل من أو يساوي التاريخ والوقت الحاليين", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var digits uint64
				var kind reflect.Kind
				var f64 float64

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				// Special handling for time.Time
				if fe.Type() == reflect.TypeOf(time.Time{}) {
					t, err = ut.T("lte-datetime", fe.Field())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("lte-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("lte-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("lte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("lte-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("lte-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("lte-items", fe.Field(), c)

				case reflect.Struct:
					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))

				default:
					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag:         "eqfield",
			translation: "يجب أن يكون {0} مساويا ل {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "يجب أن يكون {0} مساويا ل {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0} لا يمكن أن يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "يجب أن يكون {0} أكبر من {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "يجب أن يكون {0} أكبر من أو يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "يجب أن يكون {0} أصغر من {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "يجب أن يكون {0} أصغر من أو يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0} لا يمكن أن يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("gt-string", "يجب أن يكون طول {0} أكبر من {1}", false); err != nil {
					return
				}
				if err = ut.Add("gt-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("gt-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("gt-number", "يجب أن يكون {0} أكبر من {1}", false); err != nil {
					return
				}
				if err = ut.Add("gt-items", "يجب أن يحتوي {0} على أكثر من {1}", false); err != nil {
					return
				}
				if err = ut.Add("gt-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("gt-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("gt-datetime", "يجب أن يكون {0} أكبر من التاريخ والوقت الحاليين", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var digits uint64
				var kind reflect.Kind
				var f64 float64

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				// Special handling for time.Time
				if fe.Type() == reflect.TypeOf(time.Time{}) {
					t, err = ut.T("gt-datetime", fe.Field())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("gt-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("gt-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("gt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("gt-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("gt-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("gt-items", fe.Field(), c)

				case reflect.Struct:
					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))

				default:
					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag:         "alpha",
			translation: "يمكن أن يحتوي {0} على أحرف أبجدية فقط",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "يمكن أن يحتوي {0} على أحرف أبجدية رقمية فقط",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "يجب أن يكون {0} قيمة رقمية صالحة",
			override:    false,
		},
		{
			tag:         "number",
			translation: "يجب أن يكون {0} رقم صالح",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "يجب أن يكون {0} عددًا سداسيًا عشريًا صالحاً",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "يجب أن يكون {0} لون HEX صالح",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "يجب أن يكون {0} لون RGB صالح",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "يجب أن يكون {0} لون RGBA صالح",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "يجب أن يكون {0} لون HSL صالح",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "يجب أن يكون {0} لون HSLA صالح",
			override:    false,
		},
		{
			tag:         "e164",
			translation: "يجب أن يكون {0} رقم هاتف صالح بتنسيق E.164",
			override:    false,
		},
		{
			tag:         "email",
			translation: "يجب أن يكون {0} عنوان بريد إلكتروني صالح",
			override:    false,
		},
		{
			tag:         "url",
			translation: "يجب أن يكون {0} رابط إنترنت صالح",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "يجب أن يكون {0} URI صالح",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "يجب أن يكون {0} سلسلة Base64 صالحة",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "يجب أن يحتوي {0} على النص '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "يجب أن يحتوي {0} على حرف واحد على الأقل من الأحرف التالية '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "لا يمكن أن يحتوي {0} على النص '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "لا يمكن أن يحتوي {0} على أي من الأحرف التالية '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "لا يمكن أن يحتوي {0} على التالي '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "يجب أن يكون {0} رقم ISBN صالح",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "يجب أن يكون {0} رقم ISBN-10 صالح",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "يجب أن يكون {0} رقم ISBN-13 صالح",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "يجب أن يكون {0} رقم ISSN صالح",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "يجب أن يكون {0} UUID صالح",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "يجب أن يكون {0} UUID صالح من النسخة 3",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "يجب أن يكون {0} UUID صالح من النسخة 4",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "يجب أن يكون {0} UUID صالح من النسخة 5",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "يجب أن يكون {0} ULID صالح من نسخة",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "يجب أن يحتوي {0} على أحرف ascii فقط",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "يجب أن يحتوي {0} على أحرف ascii قابلة للطباعة فقط",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "يجب أن يحتوي {0} على أحرف متعددة البايت",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "يجب أن يحتوي {0} على URI صالح للبيانات",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "يجب أن يحتوي {0} على إحداثيات خط عرض صالحة",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "يجب أن يحتوي {0} على إحداثيات خط طول صالحة",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "يجب أن يكون {0} رقم SSN صالح",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "يجب أن يكون {0} عنوان IPv4 صالح",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "يجب أن يكون {0} عنوان IPv6 صالح",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "يجب أن يكون {0} عنوان IP صالح",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "يجب أن يحتوي {0} على علامة CIDR صالحة",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "يجب أن يحتوي {0} على علامة CIDR صالحة لعنوان IPv4",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "يجب أن يحتوي {0} على علامة CIDR صالحة لعنوان IPv6",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "يجب أن يكون {0} عنوان TCP صالح",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "يجب أن يكون {0} عنوان IPv4 TCP صالح",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "يجب أن يكون {0} عنوان IPv6 TCP صالح",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "يجب أن يكون {0} عنوان UDP صالح",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "يجب أن يكون {0} عنوان IPv4 UDP صالح",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "يجب أن يكون {0} عنوان IPv6 UDP صالح",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "يجب أن يكون {0} عنوان IP قابل للحل",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "يجب أن يكون {0} عنوان IP قابل للحل",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "يجب أن يكون {0} عنوان IPv6 قابل للحل",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "يجب أن يكون {0} عنوان UNIX قابل للحل",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "يجب أن يحتوي {0} على عنوان MAC صالح",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "يجب أن يحتوي {0} على قيم فريدة",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "يجب أن يكون {0} لون صالح",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "يجب أن يكون {0} واحدا من [{1}]",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "json",
			translation: "يجب أن يكون {0} نص json صالح",
			override:    false,
		},
		{
			tag:         "jwt",
			translation: "يجب أن يكون {0} نص jwt صالح",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "يجب أن يكون {0} نص حروف صغيرة",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "يجب أن يكون {0} نص حروف كبيرة",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "لا يتطابق {0} مع تنسيق {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2",
			translation: "لا يتطابق {0} مع تنسيق الرمز البريدي للبلد {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2_field",
			translation: "لا يتطابق {0} مع تنسيق الرمز البريدي للبلد في حقل {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "boolean",
			translation: "يجب أن يكون {0} قيمة منطقية صالحة",
			override:    false,
		},
		{
			tag:         "image",
			translation: "يجب أن تكون {0} صورة صالحة",
			override:    false,
		},
		{
			tag:         "cve",
			translation: "يجب أن يكون {0} معرف CVE صالح",
			override:    false,
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("gte-string", "يجب أن يكون طول {0} على الأقل {1}", false); err != nil {
					return
				}
				if err = ut.Add("gte-string-character-one", "{0} حرف", false); err != nil {
					return
				}
				if err = ut.Add("gte-string-character-other", "{0} أحرف", false); err != nil {
					return
				}
				if err = ut.Add("gte-number", "{0} يجب أن يكون {1} أو أكبر", false); err != nil {
					return
				}
				if err = ut.Add("gte-items", "يجب أن يحتوي {0} على {1} على الأقل", false); err != nil {
					return
				}
				if err = ut.Add("gte-items-item-one", "{0} عنصر", false); err != nil {
					return
				}
				if err = ut.Add("gte-items-item-other", "{0} عناصر", false); err != nil {
					return
				}
				if err = ut.Add("gte-datetime", "يجب أن يكون {0} أكبر من أو يساوي التاريخ والوقت الحاليين", false); err != nil {
					return
				}
				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				var err error
				var t string
				var digits uint64
				var kind reflect.Kind
				var f64 float64

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				// Special handling for time.Time
				if fe.Type() == reflect.TypeOf(time.Time{}) {
					t, err = ut.T("gte-datetime", fe.Field())
					goto END
				}

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				switch kind {
				case reflect.String:
					var c string
					if f64 == 1 {
						c, err = ut.T("gte-string-character-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("gte-string-character-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("gte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string
					if f64 == 1 {
						c, err = ut.T("gte-items-item-one", ut.FmtNumber(f64, digits))
					} else {
						c, err = ut.T("gte-items-item-other", ut.FmtNumber(f64, digits))
					}
					if err != nil {
						goto END
					}
					t, err = ut.T("gte-items", fe.Field(), c)

				case reflect.Struct:
					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))

				default:
					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					log.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}
				return t
			},
			override: false,
		},
		{
			tag:         "gtfield",
			translation: "يجب أن يكون {0} أكبر من {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "يجب أن يكون {0} أكبر من أو يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "يجب أن يكون {0} أصغر من {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "يجب أن يكون {0} أصغر من أو يساوي {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "fqdn",
			translation: "يجب أن يكون {0} اسم نطاق مؤهل بالكامل صالح",
			override:    false,
		},
		{
			tag:         "cron",
			translation: "يجب أن يكون {0} تعبير cron صالح",
			override:    false,
		},
		{
			tag:         "md5",
			translation: "يجب أن يكون {0} تجزئة MD5 صالحة",
			override:    false,
		},
		{
			tag:         "sha256",
			translation: "يجب أن يكون {0} تجزئة SHA256 صالحة",
			override:    false,
		},
		{
			tag:         "semver",
			translation: "يجب أن يكون {0} إصدار دلالي صالح",
			override:    false,
		},
	}

	for _, t := range translations {
		if t.customTransFunc != nil && t.customRegisFunc != nil {
			if err := v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc); err != nil {
				return err
			}
			continue
		}

		if err := v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc); err != nil {
			return err
		}
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.(error).Error()
	}
	return t
}
