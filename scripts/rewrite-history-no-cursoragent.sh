#!/bin/bash
# Переписывает историю коммитов: убирает Co-authored-by: cursoragent из сообщений.
# Запускать в клонированной копии репозитория. После выполнения: git push --force-with-lease.
set -e
if [ "${REWRITE_HISTORY_YES}" != "1" ]; then
  echo "Внимание: скрипт изменит историю. Убедитесь, что работаете с копией и нет незапушенных изменений."
  read -p "Продолжить? (y/N) " -n 1 -r
  echo
  [[ $REPLY =~ ^[yY]$ ]] || exit 1
fi

export FILTER_BRANCH_SQUELCH_WARNING=1
git filter-branch -f --msg-filter 'sed "/^Co-authored-by:.*[Cc]ursor[Aa]gent/d"' --tag-name-filter cat -- --all

echo "Готово. Проверьте историю (git log), затем: git push --force-with-lease"
echo "Удалить backup refs: git for-each-ref --format=\"%(refname)\" refs/original | xargs -n1 git update-ref -d"
