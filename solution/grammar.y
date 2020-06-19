//go:generate goyacc -o grammar.go grammar.y
%{
package solution
%}

%union{
    tok tokenType
    str string
    line int
}

%start solution

%token COMMA
%token DOT
%token EQ
%token PAREN_OPEN
%token PAREN_CLOSE
%token COMMENT
%token CRLF

%token <str> IDENTIFIER
%token <str> STRING
%token <str> BARE_STRING
%token <str> COMMENT
%token <str> COMMA

%type <str> project_type_id project_name project_type project_path project_id project_section_name project_section_type
%type <str> project_section_start project_section_key project_section_value comment lvalue rvalue

%%

solution:  lines ;

lines
    : line
    | lines CRLF line
    ;

line
    :
    | header_line
    | project_start
    | project_end
    | project_section
    | project_section_key_value_pair
    ;

header_line
        : first_line
        | version
        | comment
        ;

comment : COMMENT { $$ = $1; onComment($1) };

first_line: word words ;

words: word
        | words word ;

word : IDENTIFIER { onWord($1) }
     | COMMA
     ;

version
    : lvalue EQ rvalue { onVersion($1, $3) }
    ;

lvalue : IDENTIFIER ;

rvalue : IDENTIFIER ;

project_start
        : IDENTIFIER project_type EQ project_name COMMA project_path COMMA project_id { onProject($2, $4, $6, $8) }
        ;

project_type : PAREN_OPEN project_type_id PAREN_CLOSE { $$ = $2; };

project_type_id : STRING { $$ = $1; };

project_name : STRING { $$ = $1; };

project_path : STRING { $$ = $1; };

project_id : STRING { $$ = $1; };

project_end : IDENTIFIER ;

project_section
        : IDENTIFIER project_section_start EQ project_section_type { onSection($1, $2, $4) }
        ;

project_section_start : PAREN_OPEN project_section_name PAREN_CLOSE { $$ = $2; };

project_section_name : IDENTIFIER { $$ = $1; };

project_section_type : IDENTIFIER { $$ = $1; };

project_section_key_value_pair
        : project_section_key EQ project_section_value { onSectionItem($1, $3) }
        ;

project_section_key : BARE_STRING ;

project_section_value : BARE_STRING ;

%%