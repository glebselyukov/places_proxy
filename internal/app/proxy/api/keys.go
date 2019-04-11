package api

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/pkg/conv"
)

func prepareKey(ctx *fasthttp.RequestCtx) (preparedKey string) {
	r := ctx.Request.URI().QueryArgs()

	if r.String() == "" {
		return
	}

	keys := map[string]string{}

	r.VisitAll(func(key, value []byte) {
		k, v := strings.ToLower(conv.B2S(key)), strings.ToLower(conv.B2S(value))
		if strings.Contains(k, "[]") {
			replacedKey := strings.Replace(k, "[]", "", -1)
			if item, ok := keys[replacedKey]; ok {
				keys[replacedKey] = fmt.Sprintf("%s.%s", item, v)
				return
			}
			k = replacedKey
		}
		keys[k] = v
	})

	preparedKey = fmt.Sprintf("places:%s:%s:%s", keys["term"], keys["types"], keys["locale"])
	return
}

func queryArgs(ctx *fasthttp.RequestCtx) string {
	return ctx.Request.URI().QueryArgs().String()
}
