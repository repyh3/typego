
import chalk from "chalk";
import { Println } from "go:fmt";

Println(chalk.blue("Hello world!"));
Println(chalk.red.bold("Error!"));
Println(chalk.green.underline("Success!"));
