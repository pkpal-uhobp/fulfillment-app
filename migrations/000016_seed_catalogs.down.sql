BEGIN;

DELETE FROM cargo_place_types
WHERE name IN (
    'Паллета',
    'Коробка',
    'Мешок',
    'Ящик',
    'Бочка',
    'Биг-бэг',
    'Другое'
);

DELETE FROM product_types
WHERE name IN (
    'Одежда',
    'Обувь',
    'Косметика',
    'Электроника',
    'Игрушки',
    'Товары для дома',
    'Другое'
);

COMMIT;
