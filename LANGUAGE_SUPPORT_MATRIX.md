# Syntaxia - Матрица Поддержки Языков

**Дата:** 2025-01-16  
**Статус:** 14 языков с поддержкой тестов, линтинга и анализа ошибок

---

## 📊 Полная Матрица (15 приоритетных языков)

| Язык | Popularity | Symbols | Build | Tests | Linting | Errors | Статус | Приоритет |
|------|-----------|---------|-------|-------|---------|--------|--------|-----------|
| **Go** | ⭐⭐⭐⭐⭐ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ FULL | ✅ Done |
| **TypeScript** | ⭐⭐⭐⭐⭐ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ FULL | ✅ Done |
| **JavaScript** | ⭐⭐⭐⭐⭐ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ FULL | ✅ Done |
| **Python** | ⭐⭐⭐⭐⭐ | ✅ | ❌ | ✅ | ✅ | ✅ | 🟡 80% | 🔵 Build |
| **Java** | ⭐⭐⭐⭐⭐ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ FULL | ✅ Done |
| **Rust** | ⭐⭐⭐⭐ | ✅ | ❌ | ✅ | ✅ | ✅ | 🟡 80% | 🔵 Build |
| **C#** | ⭐⭐⭐⭐ | ✅ | ❌ | ✅ | ✅ | ✅ | 🟡 80% | 🔵 Build |
| **Kotlin** | ⭐⭐⭐ | ✅ | ❌ | ✅ | ✅ | ✅ | 🟡 80% | 🔵 Build |
| **Dart** | ⭐⭐⭐ | ✅ | ❌ | ✅ | ✅ | ✅ | 🟡 80% | 🔵 Build |
| **C/C++** | ⭐⭐⭐⭐ | ❌ | ❌ | ✅ | ✅ | ✅ | 🟡 60% | 🔵 Symbols+Build |
| **PHP** | ⭐⭐⭐ | ❌ | ❌ | ✅ | ✅ | ✅ | 🟡 60% | 🔵 Symbols+Build |
| **Swift** | ⭐⭐⭐ | ❌ | ❌ | ✅ | ✅ | ✅ | 🟡 60% | 🔵 Symbols+Build |
| **Ruby** | ⭐⭐⭐ | ❌ | ❌ | ✅ | ✅ | ✅ | 🟡 60% | 🔵 Symbols+Build |
| **Scala** | ⭐⭐ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ 0% | 🟢 Phase 6 |
| **Lua** | ⭐⭐ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ 0% | 🟢 Phase 6 |

---

## ✅ Реализованные компоненты

### Test Runners (14/14 ✅)

| Язык | Runner | Статус |
|------|--------|--------|
| Go | `go_runner.go` | ✅ |
| TypeScript | `ts_runner.go` | ✅ |
| JavaScript | `js_runner.go` | ✅ |
| Python | `python_runner.go` | ✅ |
| Java | `java_runner.go` | ✅ |
| Rust | `rust_runner.go` | ✅ |
| C# | `csharp_runner.go` | ✅ |
| Kotlin | `kotlin_runner.go` | ✅ |
| Dart | `dart_runner.go` | ✅ |
| C/C++ | `cpp_runner.go` | ✅ |
| PHP | `php_runner.go` | ✅ |
| Ruby | `ruby_runner.go` | ✅ |
| Swift | `swift_runner.go` | ✅ |

### Static Analyzers (13/13 ✅)

| Язык | Файл | Инструмент | Статус |
|------|------|------------|--------|
| Go | `staticcheck.go` | staticcheck | ✅ |
| TypeScript/JS | `eslint.go` | ESLint | ✅ |
| Python | `ruff.go` | Ruff | ✅ |
| Java | `errorprone.go` | Error Prone | ✅ |
| C/C++ | `clangtidy.go` | ClangTidy | ✅ |
| Rust | `clippy.go` | Clippy | ✅ |
| Kotlin | `ktlint.go` | ktlint | ✅ |
| C# | `dotnet_format.go` | dotnet format | ✅ |
| Dart | `dart_analyze.go` | dart analyze | ✅ |
| Swift | `swiftlint.go` | SwiftLint | ✅ |
| PHP | `phpcs.go` | PHP_CodeSniffer | ✅ |
| Ruby | `rubocop.go` | RuboCop | ✅ |

### Error Analyzers (14/14 ✅)

| Язык | Файл | Статус |
|------|------|--------|
| Go | `go_error_analyzer.go` | ✅ |
| TypeScript | `typescript_error_analyzer.go` | ✅ |
| JavaScript | `javascript_error_analyzer.go` | ✅ |
| Python | `python_error_analyzer.go` | ✅ |
| Java | `java_error_analyzer.go` | ✅ |
| Rust | `rust_error_analyzer.go` | ✅ |
| C# | `csharp_error_analyzer.go` | ✅ |
| Kotlin | `kotlin_error_analyzer.go` | ✅ |
| Dart | `dart_error_analyzer.go` | ✅ |
| C/C++ | `cpp_error_analyzer.go` | ✅ |
| PHP | `php_error_analyzer.go` | ✅ |
| Ruby | `ruby_error_analyzer.go` | ✅ |
| Swift | `swift_error_analyzer.go` | ✅ |

### Symbol Analysis (10/14)

| Язык | Статус | Примечание |
|------|--------|------------|
| Go | ✅ | tree-sitter-go |
| TypeScript | ✅ | tree-sitter-typescript |
| JavaScript | ✅ | tree-sitter-javascript |
| Vue | ✅ | tree-sitter-vue |
| Python | ✅ | tree-sitter-python |
| Java | ✅ | tree-sitter-java |
| Kotlin | ✅ | tree-sitter-kotlin |
| Rust | ✅ | tree-sitter-rust |
| C# | ✅ | tree-sitter-c-sharp |
| Dart | ✅ | tree-sitter-dart |
| C/C++ | ❌ | Нужно добавить |
| PHP | ❌ | Нужно добавить |
| Ruby | ❌ | Нужно добавить |
| Swift | ❌ | Нужно добавить |

### Build Pipeline (4/12)

| Язык | Команды | Статус |
|------|---------|--------|
| Go | `go build` | ✅ |
| TypeScript | `npm run build`, `yarn build` | ✅ |
| JavaScript | `npm run build`, `yarn build` | ✅ |
| Java | `mvn compile`, `gradle build` | ✅ |
| Python | `pip install`, `poetry install` | ❌ |
| Rust | `cargo build` | ❌ |
| C# | `dotnet build` | ❌ |
| Kotlin | `gradle build`, `kotlinc` | ❌ |
| Dart | `dart compile`, `flutter build` | ❌ |
| C/C++ | `cmake`, `make`, `g++` | ❌ |
| PHP | `composer install` | ❌ |
| Ruby | `bundle install`, `rake build` | ❌ |
| Swift | `swift build`, `xcodebuild` | ❌ |

---

## 📊 Итоговая таблица готовности

| Компонент | Готово | Всего | Процент |
|-----------|--------|-------|---------|
| Test Runners | 14 | 14 | 100% ✅ |
| Static Analyzers | 13 | 13 | 100% ✅ |
| Error Analyzers | 14 | 14 | 100% ✅ |
| Symbol Analysis | 10 | 14 | 71% 🟡 |
| Build Pipeline | 4 | 12 | 33% 🔴 |

---

## 🔧 Что нужно добавить

### Symbol Analysis (4 языка)
- [ ] C/C++ - tree-sitter-c, tree-sitter-cpp
- [ ] PHP - tree-sitter-php
- [ ] Ruby - tree-sitter-ruby
- [ ] Swift - tree-sitter-swift

### Build Pipeline (8 языков)
- [ ] Python - pip/poetry/pyproject.toml
- [ ] Rust - cargo build
- [ ] C# - dotnet build
- [ ] Kotlin - gradle/kotlinc
- [ ] Dart - dart compile/flutter build
- [ ] C/C++ - cmake/make/g++
- [ ] PHP - composer
- [ ] Ruby - bundler/rake
- [ ] Swift - swift build/xcodebuild

---

## 🔮 Phase 6: Будущие улучшения

### Scala (0% → 100%)
- [ ] Symbol analyzer (tree-sitter-scala)
- [ ] Build pipeline (sbt, gradle)
- [ ] Test runner (ScalaTest, specs2)
- [ ] Static analyzer (Scalafmt, Scalafix)
- [ ] Error analyzer

### Lua (0% → 100%)
- [ ] Symbol analyzer (tree-sitter-lua)
- [ ] Build pipeline (luarocks)
- [ ] Test runner (busted, luaunit)
- [ ] Static analyzer (luacheck)
- [ ] Error analyzer
