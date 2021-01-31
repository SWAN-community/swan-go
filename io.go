/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package swan

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func readByte(b *bytes.Buffer) (byte, error) {
	d := b.Next(1)
	if len(d) != 1 {
		return 0, fmt.Errorf("'%d' bytes incorrect for Byte", len(d))
	}
	return d[0], nil
}

func writeByte(b *bytes.Buffer, i byte) error {
	return b.WriteByte(i)
}

func readUint32(b *bytes.Buffer) (uint32, error) {
	d := b.Next(4)
	if len(d) != 4 {
		return 0, fmt.Errorf("'%d' bytes incorrect for Uint32", len(d))
	}
	return binary.LittleEndian.Uint32(d), nil
}

func writeUint32(b *bytes.Buffer, i uint32) error {
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, i)
	l, err := b.Write(v)
	if err == nil {
		if l != len(v) {
			return fmt.Errorf(
				"Mismatched lengths '%d' and '%d'",
				l,
				len(v))
		}
	}
	return err
}

func readByteArray(b *bytes.Buffer) ([]byte, error) {
	l, err := readUint32(b)
	if err != nil {
		return nil, err
	}
	return b.Next(int(l)), err
}

func writeByteArray(b *bytes.Buffer, v []byte) error {
	err := writeUint32(b, uint32(len(v)))
	if err != nil {
		return err
	}
	l, err := b.Write(v)
	if err == nil {
		if l != len(v) {
			return fmt.Errorf(
				"Mismatched lengths '%d' and '%d'",
				l,
				len(v))
		}
	}
	return err
}

func readString(b *bytes.Buffer) (string, error) {
	s, err := b.ReadBytes(0)
	if err == nil {
		return string(s[0 : len(s)-1]), err
	}
	return "", err
}

func writeString(b *bytes.Buffer, s string) error {
	l, err := b.WriteString(s)
	if err == nil {

		// Validate the number of bytes written matches the number of bytes in
		// the string.
		if l != len(s) {
			return fmt.Errorf(
				"Mismatched lengths '%d' and '%d'",
				l,
				len(s))
		}

		// Write the null terminator.
		b.WriteByte(0)
	}
	return err
}
