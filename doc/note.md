# LSP


## Position

doc: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#position

Position in a text document expressed as zero-based line and zero-based character offset. 
A position is between two characters like an ‘insert’ cursor in an editor. Special values like for example -1 to denote the end of a line are not supported.

```
interface Position {
	/**
	 * Line position in a document (zero-based).
	 */
	line: uinteger;

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	character: uinteger;
}
```

字符编码:

```
/**
 * A type indicating how positions are encoded,
 * specifically what column offsets mean.
 *
 * @since 3.17.0
 */
export type PositionEncodingKind = string;

/**
 * A set of predefined position encoding kinds.
 *
 * @since 3.17.0
 */
export namespace PositionEncodingKind {

	/**
	 * Character offsets count UTF-8 code units (e.g bytes).
	 */
	export const UTF8: PositionEncodingKind = 'utf-8';

	/**
	 * Character offsets count UTF-16 code units.
	 *
	 * This is the default and must always be supported
	 * by servers
	 */
	export const UTF16: PositionEncodingKind = 'utf-16';

	/**
	 * Character offsets count UTF-32 code units.
	 *
	 * Implementation note: these are the same as Unicode code points,
	 * so this `PositionEncodingKind` may also be used for an
	 * encoding-agnostic representation of character offsets.
	 */
	export const UTF32: PositionEncodingKind = 'utf-32';
}
```

## Range

A range in a text document expressed as (zero-based) start and end positions. A range is comparable to a selection in an editor. 
Therefore, the end position is `exclusive`. If you want to specify a range that contains a line including the line ending character(s) 
then use an end position denoting the start of the next line. For example:

```
{
    start: { line: 5, character: 23 },
    end : { line: 6, character: 0 }
}
```

```
interface Range {
	/**
	 * The range's start position.
	 */
	start: Position;

	/**
	 * The range's end position.
	 */
	end: Position;
}
```
