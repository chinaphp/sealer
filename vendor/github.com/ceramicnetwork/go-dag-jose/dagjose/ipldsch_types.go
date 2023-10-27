package dagjose

// Code generated by go-ipld-prime gengo.  DO NOT EDIT.

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

var _ datamodel.Node = nil // suppress errors when this dependency is not referenced
// Type is a struct embeding a NodePrototype/Type for every Node implementation in this package.
// One of its major uses is to start the construction of a value.
// You can use it like this:
//
// 		dagjose.Type.YourTypeName.NewBuilder().BeginMap() //...
//
// and:
//
// 		dagjose.Type.OtherTypeName.NewBuilder().AssignString("x") // ...
//
var Type typeSlab

type typeSlab struct {
	Any                     _Any__Prototype
	Any__Repr               _Any__ReprPrototype
	Base64Url               _Base64Url__Prototype
	Base64Url__Repr         _Base64Url__ReprPrototype
	Bytes                   _Bytes__Prototype
	Bytes__Repr             _Bytes__ReprPrototype
	DecodedJWE              _DecodedJWE__Prototype
	DecodedJWE__Repr        _DecodedJWE__ReprPrototype
	DecodedJWS              _DecodedJWS__Prototype
	DecodedJWS__Repr        _DecodedJWS__ReprPrototype
	DecodedRecipient        _DecodedRecipient__Prototype
	DecodedRecipient__Repr  _DecodedRecipient__ReprPrototype
	DecodedRecipients       _DecodedRecipients__Prototype
	DecodedRecipients__Repr _DecodedRecipients__ReprPrototype
	DecodedSignature        _DecodedSignature__Prototype
	DecodedSignature__Repr  _DecodedSignature__ReprPrototype
	DecodedSignatures       _DecodedSignatures__Prototype
	DecodedSignatures__Repr _DecodedSignatures__ReprPrototype
	EncodedJWE              _EncodedJWE__Prototype
	EncodedJWE__Repr        _EncodedJWE__ReprPrototype
	EncodedJWS              _EncodedJWS__Prototype
	EncodedJWS__Repr        _EncodedJWS__ReprPrototype
	EncodedRecipient        _EncodedRecipient__Prototype
	EncodedRecipient__Repr  _EncodedRecipient__ReprPrototype
	EncodedRecipients       _EncodedRecipients__Prototype
	EncodedRecipients__Repr _EncodedRecipients__ReprPrototype
	EncodedSignature        _EncodedSignature__Prototype
	EncodedSignature__Repr  _EncodedSignature__ReprPrototype
	EncodedSignatures       _EncodedSignatures__Prototype
	EncodedSignatures__Repr _EncodedSignatures__ReprPrototype
	Float                   _Float__Prototype
	Float__Repr             _Float__ReprPrototype
	Int                     _Int__Prototype
	Int__Repr               _Int__ReprPrototype
	Link                    _Link__Prototype
	Link__Repr              _Link__ReprPrototype
	List                    _List__Prototype
	List__Repr              _List__ReprPrototype
	Map                     _Map__Prototype
	Map__Repr               _Map__ReprPrototype
	Raw                     _Raw__Prototype
	Raw__Repr               _Raw__ReprPrototype
	String                  _String__Prototype
	String__Repr            _String__ReprPrototype
}

// --- type definitions follow ---

// Any matches the IPLD Schema type "Any".
// Any has union typekind, which means its data model behaviors are that of a map kind.
type Any = *_Any
type _Any struct {
	x _Any__iface
}
type _Any__iface interface {
	_Any__member()
}

func (_String) _Any__member() {}
func (_Bytes) _Any__member()  {}
func (_Int) _Any__member()    {}
func (_Float) _Any__member()  {}
func (_Map) _Any__member()    {}
func (_List) _Any__member()   {}

// Bytes matches the IPLD Schema type "Bytes".  It has bytes kind.
type Bytes = *_Bytes
type _Bytes struct{ x []byte }

// DecodedJWE matches the IPLD Schema type "DecodedJWE".  It has struct type-kind, and may be interrogated like map kind.
type DecodedJWE = *_DecodedJWE
type _DecodedJWE struct {
	aad         _Base64Url__Maybe
	ciphertext  _Base64Url
	iv          _Base64Url__Maybe
	protected   _Base64Url__Maybe
	recipients  _DecodedRecipients__Maybe
	tag         _Base64Url__Maybe
	unprotected _Any__Maybe
}

// DecodedJWS matches the IPLD Schema type "DecodedJWS".  It has struct type-kind, and may be interrogated like map kind.
type DecodedJWS = *_DecodedJWS
type _DecodedJWS struct {
	link       _Link__Maybe
	payload    _Base64Url
	signatures _DecodedSignatures__Maybe
}

// DecodedRecipient matches the IPLD Schema type "DecodedRecipient".  It has struct type-kind, and may be interrogated like map kind.
type DecodedRecipient = *_DecodedRecipient
type _DecodedRecipient struct {
	header        _Any__Maybe
	encrypted_key _Base64Url__Maybe
}

// DecodedRecipients matches the IPLD Schema type "DecodedRecipients".  It has list kind.
type DecodedRecipients = *_DecodedRecipients
type _DecodedRecipients struct {
	x []_DecodedRecipient
}

// DecodedSignature matches the IPLD Schema type "DecodedSignature".  It has struct type-kind, and may be interrogated like map kind.
type DecodedSignature = *_DecodedSignature
type _DecodedSignature struct {
	header    _Any__Maybe
	protected _Base64Url__Maybe
	signature _Base64Url
}

// DecodedSignatures matches the IPLD Schema type "DecodedSignatures".  It has list kind.
type DecodedSignatures = *_DecodedSignatures
type _DecodedSignatures struct {
	x []_DecodedSignature
}

// EncodedJWE matches the IPLD Schema type "EncodedJWE".  It has struct type-kind, and may be interrogated like map kind.
type EncodedJWE = *_EncodedJWE
type _EncodedJWE struct {
	aad         _Raw__Maybe
	ciphertext  _Raw
	iv          _Raw__Maybe
	protected   _Raw__Maybe
	recipients  _EncodedRecipients__Maybe
	tag         _Raw__Maybe
	unprotected _Any__Maybe
}

// EncodedJWS matches the IPLD Schema type "EncodedJWS".  It has struct type-kind, and may be interrogated like map kind.
type EncodedJWS = *_EncodedJWS
type _EncodedJWS struct {
	payload    _Raw
	signatures _EncodedSignatures__Maybe
}

// EncodedRecipient matches the IPLD Schema type "EncodedRecipient".  It has struct type-kind, and may be interrogated like map kind.
type EncodedRecipient = *_EncodedRecipient
type _EncodedRecipient struct {
	header        _Any__Maybe
	encrypted_key _Raw__Maybe
}

// EncodedRecipients matches the IPLD Schema type "EncodedRecipients".  It has list kind.
type EncodedRecipients = *_EncodedRecipients
type _EncodedRecipients struct {
	x []_EncodedRecipient
}

// EncodedSignature matches the IPLD Schema type "EncodedSignature".  It has struct type-kind, and may be interrogated like map kind.
type EncodedSignature = *_EncodedSignature
type _EncodedSignature struct {
	header    _Any__Maybe
	protected _Raw__Maybe
	signature _Raw
}

// EncodedSignatures matches the IPLD Schema type "EncodedSignatures".  It has list kind.
type EncodedSignatures = *_EncodedSignatures
type _EncodedSignatures struct {
	x []_EncodedSignature
}

// Float matches the IPLD Schema type "Float".  It has float kind.
type Float = *_Float
type _Float struct{ x float64 }

// Int matches the IPLD Schema type "Int".  It has int kind.
type Int = *_Int
type _Int struct{ x int64 }

// Link matches the IPLD Schema type "Link".  It has link kind.
type Link = *_Link
type _Link struct{ x datamodel.Link }

// List matches the IPLD Schema type "List".  It has list kind.
type List = *_List
type _List struct {
	x []_Any
}

// Map matches the IPLD Schema type "Map".  It has map kind.
type Map = *_Map
type _Map struct {
	m map[_String]*_Any
	t []_Map__entry
}
type _Map__entry struct {
	k _String
	v _Any
}

// String matches the IPLD Schema type "String".  It has string kind.
type String = *_String
type _String struct{ x string }
