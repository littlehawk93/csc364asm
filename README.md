# LA Tech CSC 364 Assembler

*Last updated on 2018-03-26*

This assembler was originally designed to help supplement writing assembly programs for the csc364emu emulator. Like all assemblers, this program takes in assembly source code as a plain text file and generates an executable binary that can be run by the csc364emu emulator.

Unlike most assemblers, the output file is not a true executable binary file. It is binary code formatted in the [Intel HEX](https://en.wikipedia.org/wiki/Intel_HEX) file format, specifically the I8HEX specification. The emulator reads this file to initialize the memory values for the ROM before executing.

For guidance on how to write assembly programs for the csc364emu, please visit the [CSC364EMU Github page](https://github.com/littlehawk93/csc364emu) and view the README.

To run the assembler, simply provide the input source file as an argument.

    csc364asm -f /path/to/source.asm

The assembler will create the ROM file in the same directory with the same file name as the provided source code file but with a .hex extension. To specify the output file, use the optional output flag:

    csc364asm -f /path/to/source.asm -o /output/file/path.hex

If no input file is provided, the assembler will read from stdin and output the compiled code to stdout. You can redirect stdin and stdout to effectively read and write from / to files.

    csc364asm < /path/to/source.asm > /output/file/path.hex

The assembler will provide error messages with associated line numbers when an error is encountered. 
