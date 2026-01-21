const ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
const ID_LENGTH = 7;

export function generate(): string {
    let result = "";
    for (let i = 0; i < ID_LENGTH; i++) {
        result += ALPHABET[Math.floor(Math.random() * ALPHABET.length)];
    }
    return result;
}
