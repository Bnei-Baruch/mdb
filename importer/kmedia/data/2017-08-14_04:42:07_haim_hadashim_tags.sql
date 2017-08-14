-- create tags

DO $$
DECLARE   rootid BIGINT;
  DECLARE tid    BIGINT;
BEGIN
  INSERT INTO tags (uid, pattern, description) VALUES ('IgSeiMLj', 'program-topics', 'Program topics') RETURNING id INTO rootid;
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (rootid, 'he', 'נושאי תוכנית');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (rootid, 'en', 'Program Topics');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (rootid, 'ru', 'Темы Программы');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (rootid, 'es', 'Temas del Programa');

  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'ZWtu72Y1', 'hohmat-hibur') RETURNING id INTO tid;  -- 4842
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','The Wisdom of Connection');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','חכמת החיבור');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Интегральное взаимодействие');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Sabiduría de la conexión');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'BdfLfqsa', 'zugiut') RETURNING id INTO tid;  -- 4844
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Marriage');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','זוגיות');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Супружество');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Matrimonio');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'IWXXklf7', 'arvut') RETURNING id INTO tid;  -- 4845
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Mutual responsibility');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','ערבות הדדית');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Взаимное поручительство');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Responsabilidad mutua');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'UPaIJrTw', 'herum') RETURNING id INTO tid;  -- 4846
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Emergency situation in Israel');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','מצב חרום בישראל');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Чрезвычайное положение в Израиле');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Situación de emergencia en Israel');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'TXH03g4r', 'kesef') RETURNING id INTO tid;  -- 4847
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Money');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','כסף');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Деньги');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Dinero');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'AdZ2fZcv', 'briut') RETURNING id INTO tid;  -- 4848
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Health');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','בריאות');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Здоровье');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Salud');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'tg3PuOss', 'kariera') RETURNING id INTO tid;  -- 4849
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Career, businesses and organizations');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','קריירה, עסקים וארגונים');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Карьера, предприятия и организации');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Profesión, empresas y organizaciones');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'Zi1NaTOA', 'tikshoret') RETURNING id INTO tid;  -- 4850
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Communication and social networks');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','תקשורת');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Средства связи и общения');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Comunicación y Redes Sociales');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'lsRR2yb2', 'hevra-israelit') RETURNING id INTO tid;  -- 4851
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Society of Israel');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','החברה הישראלית');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Израильское общество');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','La Sociedad de Israel');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'g3ml0jum', 'tarbut-yehudit') RETURNING id INTO tid;  -- 4852
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Jewish culture');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','תרבות יהודית');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Еврейская культура');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Cultura Judía');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'g4WzYPpT', 'horut-mishpaha') RETURNING id INTO tid;  -- 4853
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Parenting & Family');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','הורות ומשפחה');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Воспитание и Семья');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Crianza y Familia');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'riw5pzDL', 'osher') RETURNING id INTO tid;  -- 4854
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Happiness');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','אושר');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Счастье');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Felicidad');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'QZnXxBju', 'kehila') RETURNING id INTO tid;  -- 4855
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','To live in society');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','חיים קהילתיים');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Жизнь общества');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Vivir en sociedad');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, '27PtMJoy', 'ani-vehevra') RETURNING id INTO tid;  -- 4856
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Society and I');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','אני והחברה');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Я и Общество');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','La Sociedad y Yo');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'fIl5s00V', 'hagim-ve-moadim') RETURNING id INTO tid;  -- 4858
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Holidays');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','חגים ומועדים');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Праздники');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Festividades');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, '9mVGMsBY', 'kabbalah-ve-mistika') RETURNING id INTO tid;  -- 4859
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Between Kabbalah and Mysticism');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','בין קבלה למיסטיקה');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Между каббалой и мистикой');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Entre Cabalá y Misticismo');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'VyiLyjod', 'ani-ve-teva') RETURNING id INTO tid;  -- 4860
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Nature and I');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','אני והטבע');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Я и природа');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','La Naturaleza y Yo');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, '7tpk33df', 'mimshal') RETURNING id INTO tid;  -- 4863
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Leadership and management');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','הנהגה וממשל');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Руководство и управление');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Liderazgo y dirección');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'PpVfdHKx', 'gisha-haim') RETURNING id INTO tid;  -- 4868
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Approach to life');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','גישה לחיים');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Подход к жизни');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Acercamiento a la vida');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'BOXsdmd0', 'arutz-haim') RETURNING id INTO tid;  -- 6936
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','The good path in life');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','ערוץ החיים הטובים');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Доброе направление в жизни');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','El buen camino en la vida');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'iHxd7Bjj', 'bitahon') RETURNING id INTO tid;  -- 6947
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Security of Israel');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','ביטחון לישראל');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Безопасность Израиля');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Seguridad de Israel');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, 'FsKINVN5', 'megamot-olamiyot') RETURNING id INTO tid;  -- 6961
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Global tendencies');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','מגמות עולמיות');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Мировые тенденции');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Tendencias Globales');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, '0db5BBS3', 'hohmat-hakabbalah') RETURNING id INTO tid;  -- 6969
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','The Wisdom of Kabbalah');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','חכמת הקבלה');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Наука Каббала');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','La Sabiduría de la Cabalá');
  
  INSERT INTO tags (parent_id, uid, pattern) VALUES (rootid, '3cI4UcAW', 'leida') RETURNING id INTO tid;  -- 7870
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'en','Pregnancy and birth');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'he','היריון ולידה');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'ru','Беременность и роды');
  INSERT INTO tag_i18n (tag_id, language, label) VALUES (tid, 'es','Embarazo y Naciemiento');

END $$;

-- remove tags

DO $$
DECLARE rootid BIGINT;
BEGIN
  SELECT id FROM tags WHERE uid = 'IgSeiMLj'INTO rootid;

  -- delete i18ns
  DELETE FROM tag_i18n WHERE tag_id = rootid;
  DELETE FROM tag_i18n WHERE tag_id IN (SELECT id FROM tags WHERE parent_id = rootid);

  -- delete unit associations
  DELETE FROM content_units_tags WHERE tag_id = rootid;
  DELETE FROM content_units_tags WHERE tag_id IN (SELECT id
                   FROM tags
                   WHERE parent_id = rootid);

  -- delete tags themselves
  DELETE FROM tags WHERE id IN (SELECT id FROM tags WHERE parent_id = rootid);
  DELETE FROM tags WHERE id = rootid;
END $$;