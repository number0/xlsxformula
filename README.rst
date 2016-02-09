xlsxformula
==================

Excel's formula lexer and parser. It is designed to use with the following packages:

* github.com/tealeg/xlsx
* github.com/shibukawa/xlsxrange

Usage
----------

* ``xlsxformula.Tokenize(formula string) ([]*xlsxformula.Token, error)``

  Split

  .. code-block:: go

     file := xlsx.OpenFile("test.xlsx")
     sheet := file.Sheet["Sheet 1"]
     tokens, err := xlsxformula.Tokenize(sheet.Rows[1].Cells[1].Formula())

License
------------

MIT
