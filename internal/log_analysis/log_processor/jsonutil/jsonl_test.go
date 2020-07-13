package jsonutil

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"bytes"
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestAppendJoinLines(t *testing.T) {
	{
		input := []byte(`foo
bar
`)
		out := AppendJoinLines(nil, input)
		require.Equal(t, `foobar`, string(out))
	}
	{
		input := []byte(`foo bar `)
		out := AppendJoinLines(nil, input)
		require.Equal(t, `foo bar `, string(out))
	}
	{
		// inplace
		input := []byte(`foo
bar
`)
		input = AppendJoinLines(input[:0], input)
		require.Equal(t, `foobar`, string(input))
	}
}

func ExampleAppendJoinLines() {
	// Use AppendJoinLines to replace new lines in-place
	msg := []byte(`
{
  "foo": "bar"
}`)
	msg = AppendJoinLines(msg[:0], msg)
	fmt.Println(string(msg))
	// Output:{  "foo": "bar"}
}

func TestNewEncoderJSONL(t *testing.T) {
	buffer := bytes.Buffer{}
	enc := NewEncoderJSONL(&buffer, nil)
	require.Equal(t, 0, enc.NumLines())
	require.NoError(t, enc.Encode("foo"))
	require.Equal(t, 1, enc.NumLines())
	require.Equal(t, buffer.String(), `"foo"`)
	msg := jsoniter.RawMessage(`{
    "foo": "bar"
}`)
	require.NoError(t, enc.Encode(msg))
	require.Equal(t, `"foo"
{    "foo": "bar"}`, buffer.String())
	require.Equal(t, 2, enc.NumLines())
	buffer.Reset()
	enc.Reset(&buffer)
	require.NoError(t, enc.Encode(msg))
	require.Equal(t, `{    "foo": "bar"}`, buffer.String())
	require.Equal(t, 1, enc.NumLines())
}