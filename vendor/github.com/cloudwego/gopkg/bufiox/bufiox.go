// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufiox

// Reader is a buffer IO interface, which provides a user-space zero-copy method to reduce memory allocation and copy overhead.
type Reader interface {
	// Next reads the next n bytes sequentially and returns a slice `p` of length `n`,
	// otherwise returns an error if it is unable to read a buffer of n bytes.
	// The returned `p` can be a shallow copy of the original buffer.
	// Must ensure that the data in `p` is not modified before calling Release.
	//
	// Callers cannot use the returned data after calling Release.
	Next(n int) (p []byte, err error)

	// ReadBinary reads up to len(p) bytes into p. It returns the number of bytes
	// read (0 <= n <= len(p)) and any error encountered. Even if Read
	// returns n < len(p), it may use all of p as scratch space during the call.
	//
	// The bs is valid even if it is after calling Release, as it's copy read.
	ReadBinary(bs []byte) (n int, err error)

	// Peek behaves the same as Next, except that it doesn't advance the reader.
	//
	// Callers cannot use the returned data after calling Release.
	Peek(n int) (buf []byte, err error)

	// Skip skips the next n bytes sequentially, otherwise returns an error if it's unable to skip a buffer of n bytes.
	Skip(n int) (err error)

	// ReadLen returns the size that has already been read.
	// ReadBinary / Next / Skip will increase the size. When the Release function is called, ReadLen is set to 0.
	ReadLen() (n int)

	// Release will free the buffer. After release, buffer read by Next/Skip/Peek is invalid.
	// Param e is used when the buffer release depend on error.
	// Release does not close and release the underlying connection, but only releases the user mode buffer that has been read.
	Release(e error) (err error)
}

// Writer is a buffer IO interface, which provides a user-space zero-copy method to reduce memory allocation and copy overhead.
type Writer interface {
	// Malloc returns a shallow copy of the write buffer with length n,
	// otherwise returns an error if it's unable to get n bytes from the write buffer.
	// Must ensure that the data written by the user to buf can be flushed to the underlying io.Writer.
	//
	// Caller cannot write data to the returned buf after calling Flush.
	Malloc(n int) (buf []byte, err error)

	// WriteBinary writes bs to the buffer, it may be a zero copy write.
	// MUST ensure that bs is not being concurrently written before calling Flush.
	// It returns err if n < len(bs), while n is the number of bytes written.
	WriteBinary(bs []byte) (n int, err error)

	// WrittenLen returns the total length of the buffer written.
	// Malloc / WriteBinary will increase the length. When the Flush function is called, WrittenLen is set to 0.
	WrittenLen() (length int)

	// Flush writes any malloc data to the underlying io.Writer, and reset WrittenLen to zero.
	Flush() (err error)
}
