package guardrails

# Политика для проверки запрещенных путей
forbidden_paths = {
    "go.mod",
    "package.json",
    "package-lock.json",
    "yarn.lock",
    "pnpm-lock.yaml",
    ".cursor/rules/cutc.mdc"
}

# Политика для проверки временных каталогов
temp_dirs = {
    "tmp",
    "temp", 
    "cache",
    "node_modules",
    "dist",
    "build"
}

# Политика для проверки файлов с секретами
secret_files = {
    ".key",
    ".pem", 
    ".p12",
    ".pfx",
    ".env",
    ".secret"
}

# Политика для проверки бинарных файлов
binary_files = {
    ".exe",
    ".dll",
    ".so",
    ".dylib", 
    ".bin",
    ".jar",
    ".war",
    ".ear"
}

# Функция проверки запрещенных путей
is_forbidden_path(path) {
    forbidden_paths[_] = path
}

# Функция проверки временных каталогов
is_temp_dir(path) {
    startswith(path, temp_dirs[_])
}

# Функция проверки файлов с секретами
is_secret_file(path) {
    endswith(path, secret_files[_])
}

# Функция проверки бинарных файлов
is_binary_file(path) {
    endswith(path, binary_files[_])
}

# Основная функция валидации пути
validate_path(path) = result {
    result := {
        "valid": true,
        "violations": []
    }
    
    # Проверяем запрещенные пути
    is_forbidden_path(path) {
        result.valid = false
        result.violations = array.concat(result.violations, [{
            "type": "forbidden_path",
            "message": sprintf("Path %s is forbidden", [path])
        }])
    }
    
    # Проверяем временные каталоги
    is_temp_dir(path) {
        result.valid = false
        result.violations = array.concat(result.violations, [{
            "type": "temp_dir", 
            "message": sprintf("Path %s is in temporary directory", [path])
        }])
    }
    
    # Проверяем файлы с секретами
    is_secret_file(path) {
        result.valid = false
        result.violations = array.concat(result.violations, [{
            "type": "secret_file",
            "message": sprintf("Path %s contains secrets", [path])
        }])
    }
    
    # Проверяем бинарные файлы
    is_binary_file(path) {
        result.valid = false
        result.violations = array.concat(result.violations, [{
            "type": "binary_file",
            "message": sprintf("Path %s is a binary file", [path])
        }])
    }
}

# Функция проверки бюджета
validate_budget(budget_type, current, limit) = result {
    result := {
        "valid": true,
        "violations": []
    }
    
    current > limit {
        result.valid = false
        result.violations = array.concat(result.violations, [{
            "type": "budget_exceeded",
            "budget_type": budget_type,
            "current": current,
            "limit": limit,
            "message": sprintf("Budget exceeded: %v %v (limit: %v)", [current, budget_type, limit])
        }])
    }
}

# Функция проверки задачи
validate_task(task_id, files, lines_changed, ephemeral_mode) = result {
    result := {
        "valid": true,
        "violations": [],
        "budget_violations": []
    }
    
    # Проверяем каждый файл
    file_violations := [violation | 
        file := files[_]
        violation := validate_path(file)
        not violation.valid
    ]
    
    # Проверяем бюджет файлов
    file_budget := validate_budget("files", count(files), 150)
    
    # Проверяем бюджет строк
    line_budget := validate_budget("lines", lines_changed, 1500)
    
    # Собираем все нарушения
    all_violations := array.concat(file_violations, [file_budget, line_budget])
    
    # Проверяем, есть ли нарушения
    count(all_violations) > 0 {
        result.valid = false
        result.violations = all_violations
    }
    
    # В ephemeral mode разрешаем некоторые нарушения
    ephemeral_mode {
        result.valid = true
        result.ephemeral_allowed = true
    }
}

# Функция проверки конфигурации
validate_config(config) = result {
    result := {
        "valid": true,
        "warnings": []
    }
    
    # Проверяем обязательные поля
    not config.fail_closed {
        result.warnings = array.concat(result.warnings, ["fail_closed should be enabled for security"])
    }
    
    not config.enable_path_validation {
        result.warnings = array.concat(result.warnings, ["path_validation should be enabled"])
    }
    
    not config.enable_budget_tracking {
        result.warnings = array.concat(result.warnings, ["budget_tracking should be enabled"])
    }
}
