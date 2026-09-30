#!/usr/bin/env python3
"""
Раскладка файлов архива reviews для LIVECORP.

ВАЖНО: этот скрипт НЕ создаёт zip.
Он только готовит папку downloads/reviews/, из которой
архив собирается на стороне сайта (CI, вручную, build-шагом).

Что делает:
  - N однотипных логов NNN.log (шум)
  - канонические отзывы 007, 012, 404
  - секретный лог 0451 с подменённой строкой
  - README.txt с инструкцией про grep
  - фиксирует mtime у всех файлов (для воспроизводимости архива)
"""

from pathlib import Path
from datetime import datetime, timezone
import os
import random
import shutil

# ---------- настройки ----------
OUT_DIR = Path("reviews")
RANDOM_LOGS = 8000
SEED = 451

# Единая «дата инцидента» для всех файлов — фиксируем mtime.
# Выбирай любую дату; главное, чтобы 0451.log был чуть «свежее».
MTIME_BASE = datetime(2024, 3, 14, 3, 14, 15, tzinfo=timezone.utc)
MTIME_SECRET = datetime(2024, 3, 14, 4, 5, 1, tzinfo=timezone.utc)  # 04:05:01

# ---------- шаблоны (те же, что были) ----------
PRODUCTS = ["Eco-Ply 1", "Eco-Ply 1, партия Ω-6", "Eco-Ply 1, партия Ω-7"]
NOTES = [
    "Шероховатость в пределах допустимого.",
    "Прозрачность слоя соответствует ТУ.",
    "Абразивность выше нормы на 12%.",
    "Образец деформирован при намокании.",
    "Субъект отказался от повторного теста.",
    "Субъект не отказался. Это настораживает.",
    "Износ поверхности: критический.",
    "Износ поверхности: «философский».",
]
CURATOR_NOTES = [
    "Примечание куратора: субъект спокоен.",
    "Примечание куратора: субъект улыбается.",
    "Примечание куратора: улыбка зафиксирована как «лояльность бренду».",
    "Примечание куратора: запрос на морковь отклонён.",
    "Примечание куратора: запись изъята.",
]

CANONICAL = {
    "007.log": """\
ТЕСТИРОВЩИК: Кролик №7
ПРОДУКТ: Eco-Ply 1, партия Ω-6

Вердикт в протоколе: «один слой — это философия». Аудиозапись дальше нечитаема.
Примечание куратора: уши изъяты в архив образцов.
""",
    "012.log": """\
ТЕСТИРОВЩИК: Кролик №012
ПРОДУКТ: Eco-Ply 1, партия Ω-7

Шероховатость выше ожидаемой. Прозрачность слоя позволяет читать соседний протокол.
Рекомендация: выдать морковь.
Статус рекомендации: отклонено — «не входит в бюджет устойчивости».
""",
    "404.log": """\
ТЕСТИРОВЩИК: Кролик №404
ПРОДУКТ: Eco-Ply 1
ОЦЕНКА: [данные изъяты]

Поначалу было странно, но потом я просто перестал чувствовать свои уши.

Примечание куратора: субъект продолжает улыбаться.
Улыбка зафиксирована как «лояльность бренду».
""",
}

SECRET_0451 = """\
ТЕСТИРОВЩИК: Кролик №0451
ПРОДУКТ: Eco-Ply 1
ОЦЕНКА: [данные изъяты]

За компьютером 3.
""" + "\n".join(
    f"[ERROR] [RABBIT] {i:03d} — превышение порога абразивности."
    for i in range(1, 8)
) + """

Примечание куратора: субъект №0451 изъят из общего пересчёта.
Причина: см. строку 1.
"""

README = """\
LIVECORP — архив отзывов тестировщиков
Комплекс Омега, внутренний доступ

Что внутри
----------
NNN.log — отчёты по каждому субъекту. Нумерация сквозная.
Один из логов содержит аномальную запись. Мы не знаем, какой.

Как искать
----------
Логи в текстовом формате. Рекомендуем grep:

    grep -r "\\[ERROR\\] \\[RABBIT\\]" .

Найденное — сличить с внутренним реестром. Если запись
не совпадает с реестром, немедленно сообщить Старшему куратору.

Не открывайте логи вне порядка нумерации.
Не возвращайтесь тем же маршрутом.

— Отдел устойчивости
"""


def make_random_log(num: int, rng: random.Random) -> str:
    rabbit = f"Кролик №{num:03d}"
    product = rng.choice(PRODUCTS)
    note = rng.choice(NOTES)
    curator = rng.choice(CURATOR_NOTES)
    noise = "\n".join(
        f"[ERROR] [RABBIT] {rng.randint(1, 999):03d} — {rng.choice(NOTES)}"
        for _ in range(rng.randint(1, 4))
    )
    return f"""\
ТЕСТИРОВЩИК: {rabbit}
ПРОДУКТ: {product}
ОЦЕНКА: {rng.choice(['удовлетворительно', 'условно', 'отклонено', '[данные изъяты]'])}

{note}
{curator}

{noise}
"""


def set_mtime(path: Path, dt: datetime) -> None:
    """Ставит atime=mtime=dt. Работает и на Linux, и на macOS, и на Windows."""
    ts = dt.timestamp()
    os.utime(path, (ts, ts))


def build() -> None:
    rng = random.Random(SEED)

    # Чистим папку, чтобы не накапливать старьё
    if OUT_DIR.exists():
        shutil.rmtree(OUT_DIR)
    OUT_DIR.mkdir(parents=True, exist_ok=True)

    # README
    readme_path = OUT_DIR / "README.txt"
    readme_path.write_text(README, encoding="utf-8")
    set_mtime(readme_path, MTIME_BASE)

    # Канонические
    for name, text in CANONICAL.items():
        p = OUT_DIR / name
        p.write_text(text, encoding="utf-8")
        set_mtime(p, MTIME_BASE)

    # Секретный
    secret_path = OUT_DIR / "0451.log"
    secret_path.write_text(SECRET_0451, encoding="utf-8")
    set_mtime(secret_path, MTIME_SECRET)   # ← чуть «свежее» остальных

    # Шум
    taken = {7, 12, 404, 451}
    for n in range(1, RANDOM_LOGS + 1):
        if n in taken:
            continue
        p = OUT_DIR / f"{n:03d}.log"
        p.write_text(make_random_log(n, rng), encoding="utf-8")
        set_mtime(p, MTIME_BASE)

    print(f"OK: {OUT_DIR} — {len(list(OUT_DIR.iterdir()))} файлов")


if __name__ == "__main__":
    build()
