/*
 * kinds.go
 * capp-parse
 *
 * Grammar node-kind and field-name constants — the complete named surface
 * of the Objective-J grammar.
 *
 * Derived mechanically from tree-sitter-objj/src/node-types.json (generated
 * by `tree-sitter generate` from grammar.js) on 2026-06-09.  After any
 * grammar change, regenerate node-types.json and diff these constants
 * against it; the lists below are maintained by hand thereafter.
 *
 * Inclusion rules applied at derivation:
 *   - Named, concrete kinds only.  Hidden rules (underscore-prefixed) are
 *     inlined by tree-sitter and never appear as a node's kind.
 *   - Supertypes are excluded: declaration, expression, pattern,
 *     primary_expression, statement.  A node always reports its concrete
 *     kind; supertype names appear only in queries, which this project
 *     does not use.
 *   - Kinds disabled in the grammar are excluded: array, decorator,
 *     jsx_element, jsx_self_closing_element (each overridden to choice();
 *     the parser cannot produce them).
 *
 * Consumers dispatch on these constants instead of string literals, so a
 * mistyped kind is a compile error rather than a silent dispatch miss.
 *
 * Created by David Richardson on Tuesday, June 9, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package capp

// ---------------------------------------------------------------------------
// Node kinds — Objective-J
// ---------------------------------------------------------------------------

const (
	NodeKindObjjAccessorAttribute       = "objj_accessor_attribute"
	NodeKindObjjAccessorName            = "objj_accessor_name"
	NodeKindObjjAccessorsDirective      = "objj_accessors_directive"
	NodeKindObjjArrayLiteral            = "objj_array_literal"
	NodeKindObjjClassForwardDeclaration = "objj_class_forward_declaration"
	NodeKindObjjClassImplementation     = "objj_class_implementation"
	NodeKindObjjDerefExpression         = "objj_deref_expression"
	NodeKindObjjDictionaryKey           = "objj_dictionary_key"
	NodeKindObjjDictionaryLiteral       = "objj_dictionary_literal"
	NodeKindObjjDictionaryPair          = "objj_dictionary_pair"
	NodeKindObjjDictionaryValue         = "objj_dictionary_value"
	NodeKindObjjFieldDefinition         = "objj_field_definition"
	NodeKindObjjGlobalDeclaration       = "objj_global_declaration"
	NodeKindObjjImplementationMember    = "objj_implementation_member"
	NodeKindObjjImport                  = "objj_import"
	NodeKindObjjInstanceVariable        = "objj_instance_variable"
	NodeKindObjjInstanceVariables       = "objj_instance_variables"
	NodeKindObjjIvarModifier            = "objj_ivar_modifier"
	NodeKindObjjIvarModifiers           = "objj_ivar_modifiers"
	NodeKindObjjMessageExpression       = "objj_message_expression"
	NodeKindObjjMessageKeywordArgument  = "objj_message_keyword_argument"
	NodeKindObjjMessageSelector         = "objj_message_selector"
	NodeKindObjjMethodDeclaration       = "objj_method_declaration"
	NodeKindObjjMethodDefinition        = "objj_method_definition"
	NodeKindObjjMethodParameterPart     = "objj_method_parameter_part"
	NodeKindObjjMethodType              = "objj_method_type"
	NodeKindObjjMultiWordType           = "objj_multi_word_type"
	NodeKindObjjProtocolDeclaration     = "objj_protocol_declaration"
	NodeKindObjjProtocolExpression      = "objj_protocol_expression"
	NodeKindObjjProtocolReferenceList   = "objj_protocol_reference_list"
	NodeKindObjjProtocolType            = "objj_protocol_type"
	NodeKindObjjRefExpression           = "objj_ref_expression"
	NodeKindObjjSelectorExpression      = "objj_selector_expression"
	NodeKindObjjSelectorIdentifier      = "objj_selector_identifier"
	NodeKindObjjSelectorName            = "objj_selector_name"
	NodeKindObjjStringLiteral           = "objj_string_literal"
	NodeKindObjjType                    = "objj_type"
	NodeKindObjjTypedef                 = "objj_typedef"
	NodeKindObjjVisibilitySpecifier     = "objj_visibility_specifier"
)

// ---------------------------------------------------------------------------
// Node kinds — preprocessor
// ---------------------------------------------------------------------------

const (
	NodeKindPreprocCondition   = "preproc_condition"
	NodeKindPreprocConjunction = "preproc_conjunction"
	NodeKindPreprocDirective   = "preproc_directive"
	NodeKindPreprocDisjunction = "preproc_disjunction"
	NodeKindPreprocElseLine    = "preproc_else_line"
	NodeKindPreprocEndifLine   = "preproc_endif_line"
	NodeKindPreprocIfBlock     = "preproc_if_block"
	NodeKindPreprocIfLine      = "preproc_if_line"
	NodeKindPreprocIvarIfBlock = "preproc_ivar_if_block"
	NodeKindPreprocNegation    = "preproc_negation"
	NodeKindPreprocPrimary     = "preproc_primary"
)

// ---------------------------------------------------------------------------
// Node kinds — JavaScript substrate
// ---------------------------------------------------------------------------

const (
	NodeKindArguments                          = "arguments"
	NodeKindArrayPattern                       = "array_pattern"
	NodeKindArrowFunction                      = "arrow_function"
	NodeKindAssignmentExpression               = "assignment_expression"
	NodeKindAssignmentPattern                  = "assignment_pattern"
	NodeKindAugmentedAssignmentExpression      = "augmented_assignment_expression"
	NodeKindAwaitExpression                    = "await_expression"
	NodeKindBinaryExpression                   = "binary_expression"
	NodeKindBreakStatement                     = "break_statement"
	NodeKindCallExpression                     = "call_expression"
	NodeKindCatchClause                        = "catch_clause"
	NodeKindClass                              = "class"
	NodeKindClassBody                          = "class_body"
	NodeKindClassDeclaration                   = "class_declaration"
	NodeKindClassHeritage                      = "class_heritage"
	NodeKindClassStaticBlock                   = "class_static_block"
	NodeKindComment                            = "comment"
	NodeKindComputedPropertyName               = "computed_property_name"
	NodeKindContinueStatement                  = "continue_statement"
	NodeKindDebuggerStatement                  = "debugger_statement"
	NodeKindDoStatement                        = "do_statement"
	NodeKindElseClause                         = "else_clause"
	NodeKindEmptyStatement                     = "empty_statement"
	NodeKindEscapeSequence                     = "escape_sequence"
	NodeKindExportClause                       = "export_clause"
	NodeKindExportSpecifier                    = "export_specifier"
	NodeKindExportStatement                    = "export_statement"
	NodeKindExpressionStatement                = "expression_statement"
	NodeKindFalse                              = "false"
	NodeKindFieldDefinition                    = "field_definition"
	NodeKindFinallyClause                      = "finally_clause"
	NodeKindForInStatement                     = "for_in_statement"
	NodeKindForStatement                       = "for_statement"
	NodeKindFormalParameters                   = "formal_parameters"
	NodeKindFunctionDeclaration                = "function_declaration"
	NodeKindFunctionExpression                 = "function_expression"
	NodeKindGeneratorFunction                  = "generator_function"
	NodeKindGeneratorFunctionDeclaration       = "generator_function_declaration"
	NodeKindHashBangLine                       = "hash_bang_line"
	NodeKindIdentifier                         = "identifier"
	NodeKindIfStatement                        = "if_statement"
	NodeKindImport                             = "import"
	NodeKindImportAttribute                    = "import_attribute"
	NodeKindImportClause                       = "import_clause"
	NodeKindImportSpecifier                    = "import_specifier"
	NodeKindImportStatement                    = "import_statement"
	NodeKindLabeledStatement                   = "labeled_statement"
	NodeKindLexicalDeclaration                 = "lexical_declaration"
	NodeKindMemberExpression                   = "member_expression"
	NodeKindMetaProperty                       = "meta_property"
	NodeKindMethodDefinition                   = "method_definition"
	NodeKindNamedImports                       = "named_imports"
	NodeKindNamespaceExport                    = "namespace_export"
	NodeKindNamespaceImport                    = "namespace_import"
	NodeKindNativeArray                        = "native_array"
	NodeKindNewExpression                      = "new_expression"
	NodeKindNull                               = "null"
	NodeKindNumber                             = "number"
	NodeKindObject                             = "object"
	NodeKindObjectAssignmentPattern            = "object_assignment_pattern"
	NodeKindObjectPattern                      = "object_pattern"
	NodeKindOptionalChain                      = "optional_chain"
	NodeKindPair                               = "pair"
	NodeKindPairPattern                        = "pair_pattern"
	NodeKindParenthesizedExpression            = "parenthesized_expression"
	NodeKindPrivatePropertyIdentifier          = "private_property_identifier"
	NodeKindProgram                            = "program"
	NodeKindPropertyIdentifier                 = "property_identifier"
	NodeKindRegex                              = "regex"
	NodeKindRegexFlags                         = "regex_flags"
	NodeKindRegexPattern                       = "regex_pattern"
	NodeKindRestPattern                        = "rest_pattern"
	NodeKindReturnStatement                    = "return_statement"
	NodeKindSequenceExpression                 = "sequence_expression"
	NodeKindShorthandPropertyIdentifier        = "shorthand_property_identifier"
	NodeKindShorthandPropertyIdentifierPattern = "shorthand_property_identifier_pattern"
	NodeKindSpreadElement                      = "spread_element"
	NodeKindStatementBlock                     = "statement_block"
	NodeKindStatementIdentifier                = "statement_identifier"
	NodeKindString                             = "string"
	NodeKindStringFragment                     = "string_fragment"
	NodeKindSubscriptExpression                = "subscript_expression"
	NodeKindSuper                              = "super"
	NodeKindSwitchBody                         = "switch_body"
	NodeKindSwitchCase                         = "switch_case"
	NodeKindSwitchDefault                      = "switch_default"
	NodeKindSwitchStatement                    = "switch_statement"
	NodeKindSystemLibString                    = "system_lib_string"
	NodeKindTemplateString                     = "template_string"
	NodeKindTemplateSubstitution               = "template_substitution"
	NodeKindTernaryExpression                  = "ternary_expression"
	NodeKindThis                               = "this"
	NodeKindThrowStatement                     = "throw_statement"
	NodeKindTrue                               = "true"
	NodeKindTryStatement                       = "try_statement"
	NodeKindUnaryExpression                    = "unary_expression"
	NodeKindUndefined                          = "undefined"
	NodeKindUpdateExpression                   = "update_expression"
	NodeKindUsingDeclaration                   = "using_declaration"
	NodeKindVariableDeclaration                = "variable_declaration"
	NodeKindVariableDeclarator                 = "variable_declarator"
	NodeKindWhileStatement                     = "while_statement"
	NodeKindWithStatement                      = "with_statement"
	NodeKindYieldExpression                    = "yield_expression"
)

// ---------------------------------------------------------------------------
// Field names
// ---------------------------------------------------------------------------

const (
	FieldAccessors     = "accessors"
	FieldAlias         = "alias"
	FieldAlternative   = "alternative"
	FieldArgument      = "argument"
	FieldArguments     = "arguments"
	FieldBody          = "body"
	FieldCallee        = "callee"
	FieldCategory      = "category"
	FieldCondition     = "condition"
	FieldConsequence   = "consequence"
	FieldConstructor   = "constructor"
	FieldDeclaration   = "declaration"
	FieldDecorator     = "decorator"
	FieldElse          = "else"
	FieldEndif         = "endif"
	FieldFinalizer     = "finalizer"
	FieldFlags         = "flags"
	FieldFunction      = "function"
	FieldHandler       = "handler"
	FieldIf            = "if"
	FieldIncrement     = "increment"
	FieldIndex         = "index"
	FieldInitializer   = "initializer"
	FieldKey           = "key"
	FieldKeyword       = "keyword"
	FieldKind          = "kind"
	FieldLabel         = "label"
	FieldLeft          = "left"
	FieldMember        = "member"
	FieldMethodName    = "method_name"
	FieldMethodType    = "method_type"
	FieldName          = "name"
	FieldNamePart      = "name_part"
	FieldObject        = "object"
	FieldOperator      = "operator"
	FieldOptionalChain = "optional_chain"
	FieldParameter     = "parameter"
	FieldParameterName = "parameter_name"
	FieldParameterType = "parameter_type"
	FieldParameters    = "parameters"
	FieldPath          = "path"
	FieldPattern       = "pattern"
	FieldProperty      = "property"
	FieldProtocol      = "protocol"
	FieldProtocols     = "protocols"
	FieldReceiver      = "receiver"
	FieldReturnType    = "return_type"
	FieldRight         = "right"
	FieldSelector      = "selector"
	FieldSource        = "source"
	FieldSuperclass    = "superclass"
	FieldType          = "type"
	FieldValue         = "value"
)

// FieldImportPath is the original name for the objj_import path field.
// Deprecated: use FieldPath. Retained for existing consumers.
const FieldImportPath = FieldPath
