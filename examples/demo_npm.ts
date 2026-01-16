
import _ from "lodash";
import { Println } from "go:fmt";

const message = "Hello TypeGo with NPM!";
Println(_.kebabCase(message)); // Expected: hello-type-go-with-npm
Println(_.toUpper(message));   // Expected: HELLO TYPEGO WITH NPM!
