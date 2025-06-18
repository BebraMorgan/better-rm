package localization

import (
	"fmt"
	"os"
)

var messages = map[string]map[string]string{
	"en_US.UTF-8": {
		"err_remove_root":                  "Removing root directory is forbidden",
		"err_remove_trash_self":            "Removing trash directory without recovery",
		"confirm_delete_files":             "Delete %d files? (y/N)",
		"confirm_delete_file":              "Delete %s?",
		"delete_cancelled":                 "Operation cancelled by user",
		"file_deleted_verbose":             "File %s successfully moved to trash",
		"error_moving_to_trash":            "Error moving file %s to trash: %v",
		"usage_header":                     "Usage: %s [options] [files...]",
		"flag_interactive_i":               "Prompt before every removal",
		"flag_interactive_I":               "Prompt once before removing more than three files, or when removing recursively",
		"flag_verbose":                     "Explain what is being done",
		"flag_help":                        "Display this help and exit",
		"flag_version":                     "Output version information and exit",
		"flag_empty_trash":                 "Empty trash",
		"no_files_selected":                "No files selected for restoration",
		"cannot_open_trash":                "Cannot open trash",
		"could_not_find_original_path":     "Could not find original path for %s",
		"error_restoring_file":             "Error restoring file %s: %v",
		"file_restored_successfully":       "File %s restored successfully",
		"deletion_cancelled_by_user":       "Deletion cancelled by user",
		"restoration_only_in_trash":        "Can only restore files from trash",
		"unable_to_get_trash_info_path":    "Unable to get trash info path: %v",
		"unable_to_load_trash_info":        "Unable to load trash info: %v",
		"unable_to_update_trash_info":      "Unable to update trash info: %v",
		"failed_to_refresh_directory":      "Failed to refresh directory: %v",
		"error_deleting_file":              "Error deleting %s: %v",
		"error_refreshing_directory":       "Error refreshing directory: %v",
		"unable_to_determine_home_dir":     "Cannot determine home directory: %v",
		"trash_directory_created":          "Trash directory created at %s",
		"selected_file":                    "Selected: %s",
		"visual_mode_activated":            "Visual mode activated. Use ↑↓ to select multiple items.",
		"visual_mode_deactivated":          "Visual mode deactivated.",
		"could_not_find_original_path_for": "Could not find original path for",
	},

	"ru_RU.UTF-8": {
		"err_remove_root":                  "Удаление корневой директории запрещено",
		"err_remove_trash_self":            "Удаление директории корзины без возможности",
		"confirm_delete_files":             "Удалить %d файлов? (y/N)",
		"confirm_delete_file":              "Удалить %s?",
		"delete_cancelled":                 "Операция отменена пользователем",
		"file_deleted_verbose":             "Файл %s успешно перемещён в корзину",
		"error_moving_to_trash":            "Ошибка при перемещении файла %s в корзину: %v",
		"usage_header":                     "Использование: %s [опции] [файлы...]",
		"flag_interactive_i":               "Запрашивать подтверждение перед каждым удалением",
		"flag_interactive_I":               "Запрашивать подтверждение один раз перед удалением более трёх файлов или рекурсивным удалением",
		"flag_verbose":                     "Показывать подробности выполняемых действий",
		"flag_help":                        "Показать эту справку и выйти",
		"flag_version":                     "Показать информацию о версии и выйти",
		"flag_empty_trash":                 "Очистить корзину",
		"no_files_selected":                "Не выбрано ни одного файла для восстановления",
		"cannot_open_trash":                "Не удалось открыть корзину",
		"could_not_find_original_path":     "Не удалось найти оригинальный путь для %s",
		"error_restoring_file":             "Ошибка при восстановлении файла %s: %v",
		"file_restored_successfully":       "Файл %s успешно восстановлен",
		"deletion_cancelled_by_user":       "Операция отменена пользователем",
		"restoration_only_in_trash":        "Восстановление возможно только из корзины",
		"unable_to_get_trash_info_path":    "Не удалось получить путь к информации о корзине: %v",
		"unable_to_load_trash_info":        "Не удалось загрузить информацию о корзине: %v",
		"unable_to_update_trash_info":      "Не удалось обновить информацию о корзине: %v",
		"failed_to_refresh_directory":      "Не удалось обновить содержимое директории: %v",
		"error_deleting_file":              "Ошибка при удалении %s: %v",
		"error_refreshing_directory":       "Ошибка при обновлении директории: %v",
		"unable_to_determine_home_dir":     "Не удалось определить домашнюю директорию: %v",
		"trash_directory_created":          "Директория корзины создана по пути %s",
		"selected_file":                    "Выбрано: %s",
		"visual_mode_activated":            "Режим выделения активирован. Используйте ↑↓ для выбора нескольких элементов.",
		"visual_mode_deactivated":          "Режим выделения деактивирован.",
		"could_not_find_original_path_for": "Could not find original path for",
	},
}
var langCode = ""

func getLocale() string {
	if langCode == "" {
		lang := os.Getenv("LANG")
		if len(lang) >= 5 {
			langCode = lang[:5] + ".UTF-8"
		} else {
			langCode = "en_US.UTF-8"
		}
	}
	return langCode
}

func GetMessage(key string, args ...any) string {
	locale := getLocale()
	if msgs, ok := messages[locale]; ok {
		if msg, ok := msgs[key]; ok {
			return fmt.Sprintf(msg, args...)
		}
	}
	if msgs, ok := messages["en_US.UTF-8"]; ok {
		if msg, ok := msgs[key]; ok {
			return fmt.Sprintf(msg, args...)
		}
	}
	return key
}
