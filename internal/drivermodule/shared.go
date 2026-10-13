//go:build dbmod

package drivermodule

import (
	"container/list"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"os/user"
	"unicode"

	"github.com/linlay/cli-dbx/internal/conn"
	"github.com/pkujhd/goloader"
	textunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/transform"
)

// registerSharedHostTypes retains common implementations in the executable so
// the offline packer can reference them instead of copying them into modules.
// RegSymbol must run first: it initializes the loader's runtime type lookup.
// These function values retain code; registration does not call the functions.
// Keep vendor protocols out of this list so each database remains optional.
func registerSharedHostTypes(symbols map[string]uintptr) {
	goloader.RegTypes(symbols,
		tls.Dial, sql.Open, json.Unmarshal, sql.Drivers, (*sql.DB).Close,
		func(*conn.Spec) error { return nil },
		tls.X509KeyPair, x509.ParsePKCS1PrivateKey, x509.ParsePKCS8PrivateKey,
		x509.ParseECPrivateKey, x509.DecryptPEMBlock, x509.SystemCertPool,
		rand.Read, pem.Encode, pem.EncodeToMemory, binary.Write, user.Current,
		language.Parse, language.NewMatcher, message.NewPrinter,
		(*message.Printer).Sprintf, textunicode.UTF16, transform.String,
		(*json.Decoder).Decode, (*json.Decoder).Token,
		(*json.Decoder).Buffered, (*json.Decoder).More,
		(*big.Rat).FloatString, (*big.Rat).Mul, (*big.Rat).Quo,
		(*big.Rat).Set, (*big.Rat).SetInt64,
		(*big.Float).Parse, (*big.Float).Mul, (*big.Float).Quo, (*big.Float).SetPrec,
		(*list.List).PushBackList, (*list.List).PushFrontList,
		(*list.List).PushBack, (*list.List).PushFront,
		(*list.List).InsertBefore, (*list.List).InsertAfter,
		(*list.List).MoveAfter, (*list.List).MoveBefore,
		(*list.List).MoveToBack, (*list.List).MoveToFront, (*list.List).Remove,
		unicode.Categories, unicode.Scripts, unicode.Properties,
		unicode.FoldCategory, unicode.FoldScript, textunicode.UTF8, textunicode.UTF8BOM,
	)
}
