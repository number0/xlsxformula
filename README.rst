xlsxformula
==================

Excel's formula lexer and parser. It is designed to use with the following packages:

* github.com/tealeg/xlsx
* github.com/shibukawa/xlsxrange

Usage
----------

* ``xlsxformula.Tokenize(formula string) ([]*xlsxformula.Token, error)``

  Split Excel formula into Token.

  .. code-block:: go

     file := xlsx.OpenFile("test.xlsx")
     sheet := file.Sheet["Sheet 1"]
     tokens, err := xlsxformula.Tokenize(sheet.Rows[1].Cells[1].Formula())

* ``xlsxformula.Parse(formula string) ([]*xlsxformula.Node, error)``

  Split Excel formula into Token and parse structure. ``Parse()`` and ``Tokenize()`` are similar,
  but ``Parse()`` creates AST. 

  .. code-block:: go

     file := xlsx.OpenFile("test.xlsx")
     sheet := file.Sheet["Sheet 1"]
     node, err := xlsxformula.Parse(sheet.Rows[1].Cells[1].Formula())

License
------------

MIT
