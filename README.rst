xlsxformula
==================

Excel's formula lexer and parser. It is designed to use with the following packages:

* github.com/tealeg/xlsx
* github.com/shibukawa/xlsxrange

Usage
----------

* ``xlsxformula.Tokenize(formula string) ([]*xlsxformula.Token, error)``

  Split Excel formula into sequential Tokens.

  .. code-block:: go

     file := xlsx.OpenFile("test.xlsx")
     sheet := file.Sheet["Sheet 1"]
     tokens, err := xlsxformula.Tokenize(sheet.Rows[1].Cells[1].Formula())

* ``type xlsxformula.Token struct``

  * ``Type TokenType``

    It is one of the following constant values:

    * ``Number``, ``String``, ``Bool``, ``Operator``, ``LParen``, ``RParen``, ``Comma``, ``Comparator``, ``Name``, ``Range``

      Function name and named range become ``Name``.

  * ``Text string``

    Token text expression.

  * ``Line, Col int``

    Original location in formula.

* ``xlsxformula.Parse(formula string) ([]*xlsxformula.Node, error)``

  Split Excel formula into Token and parse its structure. ``Parse()`` and ``Tokenize()`` are similar,
  but ``Parse()`` creates AST. 

  .. code-block:: go

     file := xlsx.OpenFile("test.xlsx")
     sheet := file.Sheet["Sheet 1"]
     node, err := xlsxformula.Parse(sheet.Rows[1].Cells[1].Formula())

* ``type xlsxformula.Node``

  * ``Type NodeType``

    It is one of the following constant values:

    * ``Function``, ``Expression``, ``SingleToken``

  * ``Children []*xlsxformula.Node``

    * If ``NodeType`` is ``Function``, it means function's parameters.
    * If ``NodeType`` is ``Expression``, it contains other nodes (``Expression``, ``Function``, ``SingleToken``).
    * If ``NodeType`` is ``SingleToken``, it is empty.

  * ``Token *xlsxformula.Token``

    * If ``NodeType`` is ``Function``, it is ``Name`` token  as a function name.
    * If ``NodeType`` is ``Expression``, it is ``nil``.
    * If ``NodeType`` is ``SingleToken`` nodes, it is a one of ``Number``, ``String``, ``Bool``, ``Operator``, ``Comparator``, ``Name``, ``Range`` tokens.

License
------------

MIT
