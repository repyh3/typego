// MODULE: go:crypto
declare module "go:crypto" {
    export function Sha256(data: string): string;
    export function Sha512(data: string): string;
    export function HmacSha256(key: string, data: string): string;
    export function HmacSha256Verify(key: string, data: string, signature: string): boolean;
    export function RandomBytes(size: number): string;
    export function Uuid(): string;
}
// END: go:crypto
