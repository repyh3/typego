package intrinsics

import _ "embed"

//go:embed sizeof.d.ts
var sizeofTypes string

//go:embed panic.d.ts
var panicTypes string

//go:embed defer.d.ts
var deferTypes string

//go:embed scope.d.ts
var scopeTypes string

//go:embed recover.d.ts
var recoverTypes string

//go:embed pointers.d.ts
var pointerTypes string

//go:embed concurrency.d.ts
var concurrencyTypes string

//go:embed slices.d.ts
var sliceTypes string

//go:embed iota.d.ts
var iotaTypes string

//go:embed encoding.d.ts
var encodingTypes string

//go:embed buffer.d.ts
var bufferTypes string

//go:embed io.d.ts
var ioTypes string

//go:embed process.d.ts
var processTypes string

//go:embed timers.d.ts
var timerTypes string

var IntrinsicTypes = sizeofTypes + "\n" + panicTypes + "\n" + deferTypes + "\n" + scopeTypes + "\n" + recoverTypes + "\n" + pointerTypes + "\n" + concurrencyTypes + "\n" + sliceTypes + "\n" + iotaTypes + "\n" + encodingTypes + "\n" + bufferTypes + "\n" + ioTypes + "\n" + processTypes + "\n" + timerTypes
