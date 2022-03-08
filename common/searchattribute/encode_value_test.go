// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package searchattribute

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/converter"

	"go.temporal.io/server/common/payload"
)

func Test_DecodeValue_FromMetadata_Success(t *testing.T) {
	assert := assert.New(t)

	payloadStr := payload.EncodeString("qwe")
	payloadStr.Metadata["type"] = []byte("Text")
	decodedStr, err := DecodeValue(payloadStr, enumspb.INDEXED_VALUE_TYPE_UNSPECIFIED)
	assert.NoError(err)
	assert.Equal("qwe", decodedStr)

	payloadInt, err := payload.Encode(123)
	assert.NoError(err)
	payloadInt.Metadata["type"] = []byte("Int")
	decodedInt, err := DecodeValue(payloadInt, enumspb.INDEXED_VALUE_TYPE_TEXT) // MetadataType should be used anyway
	assert.NoError(err)
	assert.Equal(int64(123), decodedInt)

	payloadBool, err := payload.Encode(true)
	assert.NoError(err)
	payloadBool.Metadata["type"] = []byte("Bool")
	decodedBool, err := DecodeValue(payloadBool, enumspb.INDEXED_VALUE_TYPE_UNSPECIFIED)
	assert.NoError(err)
	assert.Equal(true, decodedBool)
}

func Test_DecodeValue_FromParameter_Success(t *testing.T) {
	assert := assert.New(t)

	payloadStr := payload.EncodeString("qwe")
	decodedStr, err := DecodeValue(payloadStr, enumspb.INDEXED_VALUE_TYPE_TEXT)
	assert.NoError(err)
	assert.Equal("qwe", decodedStr)

	payloadInt, err := payload.Encode(123)
	assert.NoError(err)
	decodedInt, err := DecodeValue(payloadInt, enumspb.INDEXED_VALUE_TYPE_INT)
	assert.NoError(err)
	assert.Equal(int64(123), decodedInt)

	payloadInt, err = payload.Encode(123)
	assert.NoError(err)
	payloadInt.Metadata["type"] = []byte("UnknownType") // should not be used because incorrect
	decodedInt, err = DecodeValue(payloadInt, enumspb.INDEXED_VALUE_TYPE_INT)
	assert.NoError(err)
	assert.Equal(int64(123), decodedInt)

	payloadBool, err := payload.Encode(true)
	assert.NoError(err)
	decodedBool, err := DecodeValue(payloadBool, enumspb.INDEXED_VALUE_TYPE_BOOL)
	assert.NoError(err)
	assert.Equal(true, decodedBool)
}

func Test_DecodeValue_Error(t *testing.T) {
	assert := assert.New(t)

	payloadStr := payload.EncodeString("qwe")
	decodedStr, err := DecodeValue(payloadStr, enumspb.INDEXED_VALUE_TYPE_UNSPECIFIED)
	assert.Error(err)
	assert.True(errors.Is(err, ErrInvalidType))
	assert.Nil(decodedStr)

	payloadInt, err := payload.Encode(123)
	assert.NoError(err)
	payloadInt.Metadata["type"] = []byte("UnknownType")
	decodedInt, err := DecodeValue(payloadInt, enumspb.INDEXED_VALUE_TYPE_UNSPECIFIED)
	assert.Error(err)
	assert.True(errors.Is(err, ErrInvalidType))
	assert.Nil(decodedInt)

	payloadInt, err = payload.Encode(123)
	assert.NoError(err)
	payloadInt.Metadata["type"] = []byte("Text")
	decodedInt, err = DecodeValue(payloadInt, enumspb.INDEXED_VALUE_TYPE_UNSPECIFIED)
	assert.Error(err)
	assert.True(errors.Is(err, converter.ErrUnableToDecode), err.Error())
	assert.Nil(decodedInt)
}

func Test_EncodeValue(t *testing.T) {
	assert := assert.New(t)

	encodedPayload, err := EncodeValue(123, enumspb.INDEXED_VALUE_TYPE_INT)
	assert.NoError(err)
	assert.Equal("123", string(encodedPayload.GetData()))
	assert.Equal("Int", string(encodedPayload.Metadata["type"]))

	encodedPayload, err = EncodeValue("qwe", enumspb.INDEXED_VALUE_TYPE_TEXT)
	assert.NoError(err)
	assert.Equal(`"qwe"`, string(encodedPayload.GetData()))
	assert.Equal("Text", string(encodedPayload.Metadata["type"]))

	var expectedEncodedRepresentation = "2022-03-07T21:27:35.986848-05:00"
	timeValue, err := time.Parse(time.RFC3339, expectedEncodedRepresentation)
	assert.NoError(err)
	encodedPayload, err = EncodeValue(timeValue, enumspb.INDEXED_VALUE_TYPE_DATETIME)
	assert.NoError(err)
	assert.Equal(`"`+expectedEncodedRepresentation+`"`, string(encodedPayload.GetData()),
		"Datetime Search Attribute is expected to be encoded in RFC 3339 format")
	assert.Equal("Datetime", string(encodedPayload.Metadata["type"]))
}
