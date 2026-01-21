/**
 * Example: External Go Module Test
 * Tests importing structs and methods from github.com/fatih/color
 */
import { Println } from "go:fmt";
import { Red, Green, Blue, Yellow } from "go:github.com/fatih/color";

Println("🎨 External Go Module Test - fatih/color");
Println("==========================================");

Red("This text should be RED");
Green("This text should be GREEN");
Blue("This text should be BLUE");
Yellow("This text should be YELLOW");

Println("");
Println("✅ External module import successful!");
Println("✅ Combined go:fmt and external module works!");
