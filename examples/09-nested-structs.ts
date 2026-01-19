/**
 * Example: Native Go Nested Struct Test
 * Tests importing structs from Go stdlib (net/url)
 */
import { Println } from "go:fmt";
import { Parse } from "go:net/url";

Println("🔗 Native Go Nested Struct Test (net/url)");
Println("==========================================");

// Parse a URL - returns a *url.URL struct with nested fields
const u = Parse("https://user:pass@example.com:8080/path?query=value#fragment");

Println("Full URL parsed!");
Println("  Scheme:", u.Scheme);
Println("  Host:", u.Host);
Println("  Path:", u.Path);
Println("  RawQuery:", u.RawQuery);
Println("  Fragment:", u.Fragment);

// Access nested User struct (Userinfo)
if (u.User) {
    Println("  User:", u.User.Username());
    // Note: Password() returns (string, bool) tuple in Go
}

// Call method on the URL struct
Println("  String():", u.String());
Println("  Hostname():", u.Hostname());
Println("  Port():", u.Port());

Println("");
Println("✅ Native stdlib nested struct access works!");
