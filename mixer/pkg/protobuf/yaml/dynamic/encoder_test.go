package dynamic

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	yaml2 "gopkg.in/yaml.v2"

	protoyaml "istio.io/istio/mixer/pkg/protobuf/yaml"
	"istio.io/istio/mixer/pkg/protobuf/yaml/testdata/all"
)

func TestEncodeVarintZeroExtend(t *testing.T) {
	for _, tst := range []struct {
		x int
		l int
	}{
		{259, 1},
		{259, 2},
		{259, 3},
		{259, 4},
		{7, 1},
		{7, 2},
		{7, 3},
		{10003432, 3},
		{10003432, 4},
		{10003432, 5},
		{10003432, 8},
	} {
		name := fmt.Sprintf("x=%v,l=%v", tst.x, tst.l)
		t.Run(name, func(tt *testing.T) {
			testEncodeVarintZeroExtend(uint64(tst.x), tst.l, tt)
		})
	}
}

func testEncodeVarintZeroExtend(x uint64, l int, tt *testing.T) {
	ba := make([]byte, 0, 8)
	ba = EncodeVarintZeroExtend(ba, x, l)

	if len(ba) < l {
		tt.Fatalf("Incorrect length. got %v, want %v", len(ba), l)
	}
	x1, n := proto.DecodeVarint(ba)
	if x1 != x {
		tt.Fatalf("Incorrect decode. got %v, want %v", x1, x)
	}
	if n != len(ba) {
		tt.Fatalf("Not all args were consumed. got %v, want %v, %v", n, len(ba), ba)
	}
}

const simple = `
str: "mystring"
i64: 56789
dbl: 123.456
b: true
enm: TWO
oth:
  str: "mystring2"
  i64: 33333
  dbl: 333.333
  b: true
  inenum: INNERTHREE
  inmsg:
    str: "myinnerstring"
    i64: 99
    dbl: 99.99
`
const sff2 = `
str: Str
dbl: 0.0021
i64: 10000203
b: true
oth:
    str: Oth.Str
mapStrStr:
    kk1: vv1
    kk2: vv2
mapI32Msg:
    200:
        str: Str
`
const sff = `
str: Str
dbl: 0.0021
i64: 10000203
b: true
`
const smm = `
mapStrStr:
    kk1: vv1
    kk2: vv2
`

type testdata struct {
	desc  string
	input string
	msg   string
}

func TestEncoder(t *testing.T) {
	fds, err := protoyaml.GetFileDescSet("../testdata/all/types.descriptor")
	if err != nil {
		t.Fatal(err)
	}
	res := protoyaml.NewResolver(fds)

	for _, td := range []testdata{
		{
			desc:  "map-only",
			msg:   ".foo.Simple",
			input: smm,
		},
		{
			desc:  "elementary",
			msg:   ".foo.Simple",
			input: sff,
		},
		{
			desc:  "full-message",
			msg:   ".foo.Simple",
			input: sff2,
		},
		{
			desc:  "full-message2",
			msg:   ".foo.Simple",
			input: simple,
		},
	} {
		t.Run(td.desc, func(tt *testing.T) {
			testMsg(tt, td.input, res)
		})
	}
}

func testMsg(t *testing.T, input string, res protoyaml.Resolver) {
	data := map[interface{}]interface{}{}
	var err error

	if err = yaml2.Unmarshal([]byte(input), data); err != nil {
		t.Fatalf("unable to unmarshal: %v\n%s", err, input)
	}

	var ba []byte

	if ba, err = yaml.YAMLToJSON([]byte(input)); err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	ff1 := foo.Simple{}
	if err = jsonpb.UnmarshalString(string(ba), &ff1); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	t.Logf("ff1 = %v", ff1)

	db := NewEncoderBuilder(".foo.Simple", res, data, nil, false)
	de, err := db.Build()

	if err != nil {
		t.Fatalf("unable to build: %v", err)
	}

	ba = make([]byte, 0, 30)
	ba, err = de.Encode(nil, ba)
	if err != nil {
		t.Fatalf("unable to encode: %v", err)
	}
	t.Logf("ba = %v", ba)

	ff2 := foo.Simple{}
	err = ff2.Unmarshal(ba)
	if err != nil {
		t.Fatalf("unable to decode: %v", err)
	}

	// confirm that codegen'd code direct unmarshal and unmarhal thru bytes yields the same result.
	if !reflect.DeepEqual(ff2, ff1) {
		t.Fatalf("got: %v, want: %v", ff2, ff1)
	}

	t.Logf("ff2 = %v", ff2)

	ba, err = ff1.Marshal()
	t.Logf("ba = %v", ba)
}
