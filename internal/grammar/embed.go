/*
 * internal/grammar/embed.go
 * capp-parse
 *
 * Created by David Richardson on Friday, April 24, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */

/*
 * EmbeddedLibrary holds the compiled Objective-J grammar dynamic library.
 * The embedded binary is the ARM Darwin build; it is used as the install
 * source by the install command and as a direct load target on ARM Darwin
 * when the library is not found in the standard search paths.
 *
 * To update: replace blobs/libtree-sitter-objj.dylib with the new build
 * and rebuild.
 */
package grammar

import _ "embed"

//go:embed blobs/libtree-sitter-objj.dylib
var EmbeddedLibrary []byte
