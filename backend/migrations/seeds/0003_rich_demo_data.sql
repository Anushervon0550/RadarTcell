-- Богатый демо-набор для RadarTcell. Перетирает технологии и каталог,
-- наполняя их реалистичными русскоязычными описаниями, картинками,
-- организациями, тегами и SDG. Изображения берутся со стабильных CDN.

BEGIN;

TRUNCATE
    technology_organizations,
    technology_tags,
    technology_sdgs,
    technologies,
    organizations,
    tags,
    sdgs,
    trends
    CASCADE;

-- ==========================================================================
-- TRENDS (сектора радара)
-- ==========================================================================
INSERT INTO trends (slug, name, description, image_url, order_index) VALUES
('ai',         'AI и машинное обучение',
                'Большие языковые модели, генеративные сети, прикладной ML в телекоме и продуктах.',
                'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=1200&q=80', 1),
('networks',   'Сети нового поколения',
                '5G, Open RAN, частные сети и эволюция в сторону 6G.',
                'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&q=80', 2),
('iot',        'Интернет вещей',
                'NB-IoT, LoRaWAN, индустриальный IoT и умный город.',
                'https://images.unsplash.com/photo-1558002038-1055907df827?w=1200&q=80', 3),
('cyber',      'Кибербезопасность',
                'Zero Trust, защита периметра оператора, SOC и обнаружение угроз.',
                'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=1200&q=80', 4),
('cloud-edge', 'Облако и Edge',
                'Cloud-native платформы, MEC, контейнеры и serverless для телекома.',
                'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80', 5),
('fintech',    'Финтех и платежи',
                'Мобильные кошельки, открытый банкинг, антифрод и BNPL для абонентов.',
                'https://images.unsplash.com/photo-1556742044-3c52d6e88c62?w=1200&q=80', 6);

-- ==========================================================================
-- SDGs
-- ==========================================================================
INSERT INTO sdgs (code, title, description, icon) VALUES
('SDG 03', 'Хорошее здоровье и благополучие',
           'Технологии, поддерживающие телемедицину и качество жизни.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-03.jpg'),
('SDG 04', 'Качественное образование',
           'Цифровые платформы для обучения и доступа к знаниям.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-04.jpg'),
('SDG 08', 'Достойная работа и экономический рост',
           'Цифровые рабочие места и автоматизация бизнес-процессов.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-08.jpg'),
('SDG 09', 'Индустриализация, инновации и инфраструктура',
           'Развитие телеком-инфраструктуры и инновационных платформ.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-09.jpg'),
('SDG 11', 'Устойчивые города и сообщества',
           'Smart City, IoT-решения для городов и транспорта.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-11.jpg'),
('SDG 13', 'Борьба с изменением климата',
           'Энергоэффективные сети и зелёные дата-центры.',
           'https://sdgs.un.org/sites/default/files/goals/E_SDG_Icons-13.jpg');

-- ==========================================================================
-- TAGS
-- ==========================================================================
INSERT INTO tags (slug, title, category, description) VALUES
('artificial-intelligence', 'Искусственный интеллект', 'Domain',     'Алгоритмы, имитирующие интеллект.'),
('ml',                      'Машинное обучение',        'Domain',     'Обучение моделей на данных.'),
('nlp',                     'Обработка языка',          'Domain',     'NLP и LLM.'),
('computer-vision',         'Компьютерное зрение',       'Domain',     'Анализ изображений и видео.'),
('telecom',                 'Телеком',                   'Industry',   'Операторы связи.'),
('5g',                      '5G',                        'Tech',       'Пятое поколение мобильной связи.'),
('iot',                     'IoT',                       'Tech',       'Связанные устройства.'),
('edge',                    'Edge Computing',            'Tech',       'Распределённые вычисления.'),
('cloud',                   'Cloud',                     'Tech',       'Облачные платформы.'),
('security',                'Безопасность',              'Domain',     'Защита данных и инфраструктуры.'),
('blockchain',              'Блокчейн',                  'Tech',       'Распределённый реестр.'),
('analytics',               'Аналитика',                 'Domain',     'BI и большие данные.'),
('mobile-money',            'Мобильные финансы',          'Industry',   'Платежи через мобильное устройство.'),
('automation',              'Автоматизация',             'Domain',     'Снижение ручного труда.');

-- ==========================================================================
-- ORGANIZATIONS
-- ==========================================================================
INSERT INTO organizations (slug, name, logo_url, description, website, headquarters) VALUES
('tcell',      'Tcell',
               'https://logo.clearbit.com/tcell.tj',
               'Ведущий мобильный оператор Таджикистана. Развивает 4G/5G, цифровые сервисы и финтех.',
               'https://tcell.tj', 'Душанбе, Таджикистан'),
('openai',     'OpenAI',
               'https://logo.clearbit.com/openai.com',
               'Разработчик GPT, DALL·E и инструментов для прикладного ИИ.',
               'https://openai.com', 'Сан-Франциско, США'),
('nvidia',     'NVIDIA',
               'https://logo.clearbit.com/nvidia.com',
               'Производитель GPU и AI-платформ для дата-центров и периферии.',
               'https://nvidia.com', 'Санта-Клара, США'),
('ericsson',   'Ericsson',
               'https://logo.clearbit.com/ericsson.com',
               'Глобальный поставщик телеком-решений: радиосети, 5G, BSS/OSS.',
               'https://ericsson.com', 'Стокгольм, Швеция'),
('huawei',     'Huawei',
               'https://logo.clearbit.com/huawei.com',
               'Производитель сетевого оборудования и облачных платформ.',
               'https://huawei.com', 'Шэньчжэнь, Китай'),
('cisco',      'Cisco',
               'https://logo.clearbit.com/cisco.com',
               'Сетевые решения, безопасность и совместная работа.',
               'https://cisco.com', 'Сан-Хосе, США'),
('aws',        'Amazon Web Services',
               'https://logo.clearbit.com/aws.amazon.com',
               'Крупнейший облачный провайдер: вычисления, хранение, ML, IoT.',
               'https://aws.amazon.com', 'Сиэтл, США'),
('mts',        'МТС',
               'https://logo.clearbit.com/mts.ru',
               'Российский оператор связи и цифровая экосистема.',
               'https://mts.ru', 'Москва, Россия'),
('palo-alto',  'Palo Alto Networks',
               'https://logo.clearbit.com/paloaltonetworks.com',
               'Решения по сетевой безопасности и SASE.',
               'https://paloaltonetworks.com', 'Санта-Клара, США'),
('snowflake',  'Snowflake',
               'https://logo.clearbit.com/snowflake.com',
               'Облачная платформа для данных и аналитики.',
               'https://snowflake.com', 'Бозман, США');

-- ==========================================================================
-- TECHNOLOGIES
-- ==========================================================================
WITH
    t_ai    AS (SELECT id FROM trends WHERE slug = 'ai'),
    t_net   AS (SELECT id FROM trends WHERE slug = 'networks'),
    t_iot   AS (SELECT id FROM trends WHERE slug = 'iot'),
    t_cyber AS (SELECT id FROM trends WHERE slug = 'cyber'),
    t_cloud AS (SELECT id FROM trends WHERE slug = 'cloud-edge'),
    t_fin   AS (SELECT id FROM trends WHERE slug = 'fintech')

INSERT INTO technologies (
    slug, list_index, name,
    description_short, description_full,
    readiness_level,
    custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4,
    image_url, source_link,
    trend_id
) VALUES
-- ===== AI =====
('edge-llm', 1, 'Edge LLM',
 'Большие языковые модели, работающие прямо на устройстве пользователя.',
 'Edge LLM — это компактные языковые модели (1–7B параметров), оптимизированные под запуск на смартфонах, ноутбуках и edge-серверах. Применяются для офлайн-ассистентов, обработки чувствительных данных без отправки в облако и для снижения задержек. Поддерживают инструктаж, RAG и tool-use в ограниченном по памяти окружении. Стадия: ранний продукт, активное развитие.',
 6, 0.72, 0.45, 0.80, 0.28,
 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=1200&q=80',
 'https://huggingface.co/blog/llm-on-device',
 (SELECT id FROM t_ai)),

('rag-platform', 2, 'RAG Platform',
 'Платформа поиска и генерации ответов по корпоративным базам знаний.',
 'Retrieval-Augmented Generation объединяет векторный поиск с LLM, обеспечивая ответы на основе актуальных документов компании. Снимает галлюцинации модели, прозрачно ссылается на источники, поддерживает инкрементальное обновление индекса и контроль доступа. Активно используется в поддержке клиентов, юридическом анализе и онбординге сотрудников.',
 7, 0.78, 0.62, 0.71, 0.40,
 'https://images.unsplash.com/photo-1620712943543-bcc4688e7485?w=1200&q=80',
 'https://www.pinecone.io/learn/retrieval-augmented-generation/',
 (SELECT id FROM t_ai)),

('fraud-ml', 3, 'Antifraud ML',
 'Модели машинного обучения для выявления мошеннических транзакций в реальном времени.',
 'Решение строит профили нормального поведения абонентов и счетов, выявляя аномалии в платежах и звонках. Использует градиентный бустинг, графовые нейросети и потоковую обработку (Kafka + Flink). Снижает потери от мошенничества на 40–60% и ложные срабатывания на 30%.',
 8, 0.82, 0.50, 0.45, 0.78,
 'https://images.unsplash.com/photo-1563013544-824ae1b704d3?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Credit_card_fraud',
 (SELECT id FROM t_ai)),

('voice-analytics', 4, 'Voice Analytics',
 'Аналитика речи в контакт-центре: эмоции, темы, скрипты.',
 'Распознавание речи (ASR) + NLU выделяют темы обращений, оценивают тональность абонента и соблюдение скриптов оператором. Помогает улучшить NPS, обучать сотрудников, выявлять отток. Поддерживает русский, таджикский и английский языки. Стадия: производственная.',
 7, 0.65, 0.55, 0.70, 0.48,
 'https://images.unsplash.com/photo-1589254065878-42c9da997008?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Speech_analytics',
 (SELECT id FROM t_ai)),

('recommendation-ai', 5, 'Recommendation Engine',
 'Персональные рекомендации тарифов, контента и допуслуг.',
 'Двухуровневая модель (отбор кандидатов + ранжирование) на основе матричной факторизации, последовательных трансформеров и контекстных бандитов. Увеличивает ARPU за счёт точечных предложений в приложении и push-каналах. Интегрируется с CRM и Big Data.',
 7, 0.74, 0.68, 0.72, 0.42,
 'https://images.unsplash.com/photo-1611162617213-7d7a39e9b1d7?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Recommender_system',
 (SELECT id FROM t_ai)),

('aiops', 6, 'AIOps',
 'AI для эксплуатации сетей: предиктивное обнаружение и автозакрытие инцидентов.',
 'AIOps собирает метрики, логи и события сети, обнаруживает аномалии и кореллирует инциденты. Снижает MTTR и количество ложных тикетов. Часто реализуется как платформа поверх ELK + ClickHouse + ML-сервиса. Стадия: пилоты у крупных операторов.',
 5, 0.55, 0.62, 0.50, 0.45,
 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=1200&q=80',
 'https://www.gartner.com/en/information-technology/glossary/aiops-platform',
 (SELECT id FROM t_ai)),

('churn-prediction', 7, 'Churn Prediction',
 'Прогноз оттока абонентов и сегментация для удержания.',
 'Модель на основе градиентного бустинга и исторических CDR/биллинга предсказывает вероятность ухода в ближайшие 30/60/90 дней. Сегментация по причинам оттока запускает таргетированные акции. ROI до 5x за счёт удержания high-value абонентов.',
 7, 0.60, 0.66, 0.55, 0.52,
 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Customer_attrition',
 (SELECT id FROM t_ai)),

('computer-vision-qa', 8, 'CV для контроля качества',
 'Компьютерное зрение для контроля монтажа базовых станций по фото.',
 'Модель проверяет фото с объектов: правильность маркировки, заземление, расположение антенн. Запускается в мобильном приложении монтажников. Сокращает повторные выезды на 25%. Стадия: пилот.',
 4, 0.42, 0.48, 0.62, 0.35,
 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Computer_vision',
 (SELECT id FROM t_ai)),

-- ===== Networks =====
('open-ran', 9, 'Open RAN',
 'Открытая радиоподсистема с мульти-вендорной интеграцией.',
 'Open RAN разделяет аппаратную и программную части базовой станции (RU/DU/CU), позволяя смешивать оборудование разных вендоров. Снижает CAPEX и vendor lock-in. Активно тестируется в Европе и Японии. В Таджикистане — ранние пилоты.',
 5, 0.50, 0.62, 0.30, 0.40,
 'https://images.unsplash.com/photo-1592609931095-54a2168ae893?w=1200&q=80',
 'https://www.o-ran.org/',
 (SELECT id FROM t_net)),

('network-slicing', 10, 'Network Slicing',
 'Логические подсети 5G под разные SLA: IoT, eMBB, URLLC.',
 'Slicing позволяет на одной физической сети создавать виртуальные подсети с разными гарантиями по задержке, пропускной способности и надёжности. Применяется для индустриальных сценариев, AR/VR, автономного транспорта. Стадия: коммерческое внедрение в отдельных рынках.',
 6, 0.45, 0.72, 0.62, 0.50,
 'https://images.unsplash.com/photo-1581090700227-1e37b190418e?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Network_slicing',
 (SELECT id FROM t_net)),

('private-5g', 11, 'Private 5G',
 'Частные сети 5G для предприятий: производство, порты, аэропорты.',
 'Изолированная 5G-сеть на территории заказчика обеспечивает гарантированные SLA для роботов, AGV, видеоаналитики. Спектр выделяется регулятором или операторской частной полосой. ROI достигается за счёт автоматизации производственных процессов.',
 6, 0.50, 0.66, 0.58, 0.55,
 'https://images.unsplash.com/photo-1581092334651-ddf26d9a09d0?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Private_LTE',
 (SELECT id FROM t_net)),

('mec', 12, 'Multi-access Edge Computing',
 'Перенос вычислений на границу сети для минимальной задержки.',
 'MEC размещает приложения на серверах рядом с базовыми станциями. Используется для облачных игр, AR, видеоаналитики, V2X. Стандартизирован ETSI. Активные внедрения у Telefónica, Verizon, NTT.',
 6, 0.52, 0.60, 0.66, 0.49,
 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80',
 'https://www.etsi.org/technologies/multi-access-edge-computing',
 (SELECT id FROM t_net)),

('wifi-7', 13, 'Wi-Fi 7',
 'Новый стандарт IEEE 802.11be: пиковые скорости до 46 Гбит/с.',
 'Wi-Fi 7 (802.11be) добавляет MLO (Multi-Link Operation), 320 МГц каналы и 4096-QAM. Значительно повышает throughput и стабильность для AR/VR и cloud-gaming. Сертификация Wi-Fi Alliance стартовала в 2024.',
 7, 0.62, 0.55, 0.68, 0.45,
 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&q=80',
 'https://www.wi-fi.org/discover-wi-fi/wi-fi-certified-7',
 (SELECT id FROM t_net)),

('6g-research', 14, '6G Research',
 'Исследовательские работы по сетям шестого поколения.',
 'Терагерцовый диапазон, интеграция связи и сенсинга, AI-native сети, скорости в десятки Гбит/с и задержки в микросекунды. Коммерческий запуск ожидается около 2030 года. Активные программы: Hexa-X (EU), Next G Alliance (US), IMT-2030 (3GPP).',
 3, 0.30, 0.35, 0.25, 0.20,
 'https://images.unsplash.com/photo-1551434678-e076c223a692?w=1200&q=80',
 'https://www.itu.int/en/ITU-T/focusgroups/net2030/',
 (SELECT id FROM t_net)),

-- ===== IoT =====
('nb-iot', 15, 'NB-IoT',
 'Узкополосная сотовая сеть для миллионов устройств с долгим временем работы.',
 'NB-IoT обеспечивает покрытие в подвалах и счётчиках, экономит батарею (10+ лет) и поддерживает огромную плотность устройств. Применяется в ЖКХ, логистике, сельском хозяйстве. Поддерживается в большинстве 4G-сетей.',
 8, 0.70, 0.40, 0.55, 0.62,
 'https://images.unsplash.com/photo-1558002038-1055907df827?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Narrowband_IoT',
 (SELECT id FROM t_iot)),

('lorawan', 16, 'LoRaWAN',
 'Открытый протокол LPWAN для нелицензируемого спектра.',
 'LoRaWAN обеспечивает дальность до 15 км в сельской местности при минимальном энергопотреблении. Используется для умного города, охраны окружающей среды, агро. Развёртывается оператором или частной сетью.',
 8, 0.66, 0.42, 0.50, 0.58,
 'https://images.unsplash.com/photo-1565536421951-1bdaa4f53e94?w=1200&q=80',
 'https://lora-alliance.org/',
 (SELECT id FROM t_iot)),

('smart-meters', 17, 'Smart Metering',
 'Умные счётчики электроэнергии и воды.',
 'Дистанционный сбор показаний, обнаружение утечек и хищений, интеграция с биллингом. Часто строится на NB-IoT или PLC. Снижает операционные расходы коммунальных служб на 15–25%.',
 7, 0.60, 0.45, 0.52, 0.65,
 'https://images.unsplash.com/photo-1581090700227-1e37b190418e?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Smart_meter',
 (SELECT id FROM t_iot)),

('connected-vehicles', 18, 'Connected Vehicles',
 'Подключённые автомобили: телематика, V2X, OTA.',
 'Поддержка eCall, удалённой диагностики, OTA-обновлений ПО, V2X-коммуникации. Совместно с 5G и MEC обеспечивает основу для автономного транспорта.',
 6, 0.55, 0.58, 0.50, 0.46,
 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Connected_car',
 (SELECT id FROM t_iot)),

('smart-city', 19, 'Smart City Platform',
 'Платформа цифрового города: освещение, транспорт, безопасность.',
 'Объединяет данные с тысяч датчиков, камер и сервисов, предоставляя единый Dashboard городским службам. Включает аналитику движения, мониторинг качества воздуха, предиктивное обслуживание инфраструктуры.',
 6, 0.58, 0.62, 0.66, 0.50,
 'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Smart_city',
 (SELECT id FROM t_iot)),

-- ===== Cybersecurity =====
('zero-trust', 20, 'Zero Trust Network',
 'Архитектура «нулевого доверия» для защиты корпоративных ресурсов.',
 'Zero Trust подразумевает аутентификацию каждого запроса, минимально необходимые привилегии, микросегментацию. Реализуется через ZTNA, IAM, MFA, EDR. Замещает классический периметровый подход с VPN.',
 7, 0.68, 0.60, 0.72, 0.55,
 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=1200&q=80',
 'https://www.nist.gov/publications/zero-trust-architecture',
 (SELECT id FROM t_cyber)),

('soc-automation', 21, 'SOC Automation (SOAR)',
 'Автоматизация security operations center: playbooks и оркестрация.',
 'SOAR ускоряет реагирование на инциденты за счёт готовых сценариев: блокировка IP, изоляция хоста, обогащение тикета. Снижает MTTR в 5–10 раз. Часто интегрируется с SIEM, EDR и threat intel.',
 7, 0.64, 0.58, 0.66, 0.52,
 'https://images.unsplash.com/photo-1614064641938-3bbee52942c7?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Security_orchestration',
 (SELECT id FROM t_cyber)),

('ddos-protection', 22, 'DDoS Protection',
 'Защита от распределённых атак на уровне оператора.',
 'Скраббинг-центры, BGP-flowspec, ML-обнаружение аномалий. Защита веб-сервисов, игровых платформ и инфраструктуры оператора. Поддержка Anycast и автоматическая активация.',
 8, 0.74, 0.50, 0.60, 0.68,
 'https://images.unsplash.com/photo-1563206767-5b18f218e8de?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Denial-of-service_attack',
 (SELECT id FROM t_cyber)),

('post-quantum-crypto', 23, 'Post-Quantum Crypto',
 'Криптография, устойчивая к квантовым атакам.',
 'NIST PQC: ML-KEM (Kyber), ML-DSA (Dilithium), SLH-DSA. Подготовка к замене RSA/ECC до того, как квантовые компьютеры станут практическими. Стадия: стандартизация и пилоты.',
 4, 0.40, 0.45, 0.30, 0.25,
 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=1200&q=80',
 'https://csrc.nist.gov/projects/post-quantum-cryptography',
 (SELECT id FROM t_cyber)),

-- ===== Cloud / Edge =====
('kubernetes-telco', 24, 'Kubernetes для телекома',
 'Cloud-native сетевые функции (CNF) на Kubernetes.',
 'Запуск 5G Core, IMS, BSS на K8s обеспечивает масштабирование, GitOps, быстрые обновления. Используется CNCF-стек: Helm, Argo CD, Istio, Prometheus. Совместная работа с CNCF Telecom User Group.',
 7, 0.66, 0.62, 0.70, 0.50,
 'https://images.unsplash.com/photo-1605379399642-870262d3d051?w=1200&q=80',
 'https://www.cncf.io/projects/kubernetes/',
 (SELECT id FROM t_cloud)),

('serverless-billing', 25, 'Serverless Billing',
 'Биллинг на serverless-функциях с pay-per-use.',
 'Функциональная архитектура (AWS Lambda / Knative) обрабатывает пиковые нагрузки без overprovisioning. События тарификации поступают через Kafka, агрегация в ClickHouse, отчёты в S3. Стадия: ранние внедрения.',
 5, 0.48, 0.55, 0.58, 0.42,
 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Serverless_computing',
 (SELECT id FROM t_cloud)),

('observability', 26, 'Observability Stack',
 'Сквозное наблюдение: метрики, логи, трейсы (OpenTelemetry).',
 'Единый стек на основе OpenTelemetry, Prometheus, Loki, Tempo, Grafana. Корреляция между сигналами, распределённый tracing, SLO-мониторинг. Снижает MTTD и MTTR в продакшене.',
 8, 0.76, 0.58, 0.74, 0.55,
 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=1200&q=80',
 'https://opentelemetry.io/',
 (SELECT id FROM t_cloud)),

-- ===== Fintech =====
('mobile-wallet', 27, 'Mobile Wallet',
 'Кошелёк оператора: переводы, оплата ЖКХ, QR-платежи.',
 'Электронный кошелёк интегрирован с биллингом и банками. Поддержка KYC, лимитов, программы лояльности, QR Pay. Конкурирует с банковскими приложениями. Высокий потенциал в развивающихся рынках.',
 8, 0.78, 0.55, 0.65, 0.72,
 'https://images.unsplash.com/photo-1556742044-3c52d6e88c62?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Mobile_payment',
 (SELECT id FROM t_fin)),

('open-banking', 28, 'Open Banking',
 'Открытые API для доступа к банковским данным с согласия клиента.',
 'PSD2/UK Open Banking стандартизировали доступ к счетам и инициации платежей. Оператор может выступать TPP, агрегируя счета пользователя в едином интерфейсе. Стадия: коммерческие продукты в EU/UK, пилоты в СНГ.',
 6, 0.58, 0.62, 0.55, 0.50,
 'https://images.unsplash.com/photo-1601597111158-2fceff292cdc?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Open_banking',
 (SELECT id FROM t_fin)),

('bnpl', 29, 'BNPL для абонентов',
 'Покупка устройств в рассрочку с интеграцией в биллинг.',
 'Buy Now Pay Later для смартфонов и аксессуаров, рассрочка 3–24 месяца на лицевом счёте абонента. Скоринг на основе истории платежей и behavioral-данных. Низкий уровень дефолтов за счёт автосписания.',
 7, 0.65, 0.60, 0.58, 0.62,
 'https://images.unsplash.com/photo-1563013544-824ae1b704d3?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Buy_now,_pay_later',
 (SELECT id FROM t_fin)),

('crypto-custody', 30, 'Crypto Custody',
 'Сервис безопасного хранения цифровых активов.',
 'Хранение ключей через MPC и HSM, мульти-подпись, страхование. Подготовка к интеграции CBDC и регулируемых стейблкоинов. Стадия: ранний рынок, активная регуляторная неопределённость.',
 4, 0.38, 0.42, 0.32, 0.28,
 'https://images.unsplash.com/photo-1518546305927-5a555bb7020d?w=1200&q=80',
 'https://en.wikipedia.org/wiki/Cryptocurrency_wallet',
 (SELECT id FROM t_fin));

-- ==========================================================================
-- LINKS: technology -> tags
-- ==========================================================================
INSERT INTO technology_tags (technology_id, tag_id)
SELECT t.id, tag.id FROM technologies t JOIN tags tag ON
       (t.slug = 'edge-llm'              AND tag.slug IN ('artificial-intelligence','ml','nlp','edge'))
    OR (t.slug = 'rag-platform'          AND tag.slug IN ('artificial-intelligence','ml','nlp','analytics'))
    OR (t.slug = 'fraud-ml'              AND tag.slug IN ('ml','analytics','security','automation'))
    OR (t.slug = 'voice-analytics'       AND tag.slug IN ('artificial-intelligence','nlp','analytics','telecom'))
    OR (t.slug = 'recommendation-ai'     AND tag.slug IN ('artificial-intelligence','ml','analytics'))
    OR (t.slug = 'aiops'                 AND tag.slug IN ('ml','automation','telecom','analytics'))
    OR (t.slug = 'churn-prediction'      AND tag.slug IN ('ml','analytics','telecom'))
    OR (t.slug = 'computer-vision-qa'    AND tag.slug IN ('artificial-intelligence','computer-vision','automation'))
    OR (t.slug = 'open-ran'              AND tag.slug IN ('telecom','5g','automation'))
    OR (t.slug = 'network-slicing'       AND tag.slug IN ('telecom','5g'))
    OR (t.slug = 'private-5g'            AND tag.slug IN ('telecom','5g','iot'))
    OR (t.slug = 'mec'                   AND tag.slug IN ('telecom','edge','cloud'))
    OR (t.slug = 'wifi-7'                AND tag.slug IN ('telecom'))
    OR (t.slug = '6g-research'           AND tag.slug IN ('telecom'))
    OR (t.slug = 'nb-iot'                AND tag.slug IN ('iot','telecom'))
    OR (t.slug = 'lorawan'               AND tag.slug IN ('iot'))
    OR (t.slug = 'smart-meters'          AND tag.slug IN ('iot','automation'))
    OR (t.slug = 'connected-vehicles'    AND tag.slug IN ('iot','5g','telecom'))
    OR (t.slug = 'smart-city'            AND tag.slug IN ('iot','analytics','automation'))
    OR (t.slug = 'zero-trust'            AND tag.slug IN ('security','cloud'))
    OR (t.slug = 'soc-automation'        AND tag.slug IN ('security','automation','ml'))
    OR (t.slug = 'ddos-protection'       AND tag.slug IN ('security','telecom'))
    OR (t.slug = 'post-quantum-crypto'   AND tag.slug IN ('security'))
    OR (t.slug = 'kubernetes-telco'      AND tag.slug IN ('cloud','telecom','automation'))
    OR (t.slug = 'serverless-billing'    AND tag.slug IN ('cloud','telecom','automation'))
    OR (t.slug = 'observability'         AND tag.slug IN ('cloud','analytics','automation'))
    OR (t.slug = 'mobile-wallet'         AND tag.slug IN ('mobile-money','telecom'))
    OR (t.slug = 'open-banking'          AND tag.slug IN ('mobile-money','analytics'))
    OR (t.slug = 'bnpl'                  AND tag.slug IN ('mobile-money','ml'))
    OR (t.slug = 'crypto-custody'        AND tag.slug IN ('blockchain','security','mobile-money'))
ON CONFLICT DO NOTHING;

-- ==========================================================================
-- LINKS: technology -> sdgs
-- ==========================================================================
INSERT INTO technology_sdgs (technology_id, sdg_id)
SELECT t.id, s.id FROM technologies t JOIN sdgs s ON
       (t.slug IN ('edge-llm','rag-platform','aiops','observability','kubernetes-telco','serverless-billing','open-ran','network-slicing','private-5g','mec','wifi-7','6g-research','nb-iot','lorawan') AND s.code = 'SDG 09')
    OR (t.slug IN ('voice-analytics','recommendation-ai','rag-platform') AND s.code = 'SDG 04')
    OR (t.slug IN ('connected-vehicles','smart-city','smart-meters') AND s.code = 'SDG 11')
    OR (t.slug IN ('zero-trust','soc-automation','ddos-protection','post-quantum-crypto','fraud-ml') AND s.code = 'SDG 09')
    OR (t.slug IN ('mobile-wallet','open-banking','bnpl','churn-prediction') AND s.code = 'SDG 08')
    OR (t.slug IN ('smart-meters','smart-city','observability') AND s.code = 'SDG 13')
    OR (t.slug IN ('voice-analytics','computer-vision-qa') AND s.code = 'SDG 03')
ON CONFLICT DO NOTHING;

-- ==========================================================================
-- LINKS: technology -> organizations
-- ==========================================================================
INSERT INTO technology_organizations (technology_id, organization_id)
SELECT t.id, o.id FROM technologies t JOIN organizations o ON
       (t.slug IN ('edge-llm','rag-platform','voice-analytics','recommendation-ai','computer-vision-qa') AND o.slug = 'openai')
    OR (t.slug IN ('edge-llm','recommendation-ai','aiops','computer-vision-qa') AND o.slug = 'nvidia')
    OR (t.slug IN ('open-ran','network-slicing','private-5g','mec','wifi-7','6g-research','aiops') AND o.slug = 'ericsson')
    OR (t.slug IN ('open-ran','network-slicing','private-5g','5g','wifi-7','nb-iot','lorawan','smart-meters','smart-city','connected-vehicles','6g-research') AND o.slug = 'huawei')
    OR (t.slug IN ('zero-trust','soc-automation','ddos-protection','wifi-7','observability') AND o.slug = 'cisco')
    OR (t.slug IN ('serverless-billing','observability','kubernetes-telco','rag-platform','recommendation-ai') AND o.slug = 'aws')
    OR (t.slug IN ('mobile-wallet','open-banking','bnpl','fraud-ml','churn-prediction','voice-analytics','recommendation-ai','aiops','smart-city','observability','open-ran','network-slicing','private-5g','mec','nb-iot','smart-meters','connected-vehicles','zero-trust','soc-automation','ddos-protection','kubernetes-telco') AND o.slug = 'tcell')
    OR (t.slug IN ('fraud-ml','open-banking','bnpl','mobile-wallet','aiops','observability','churn-prediction') AND o.slug = 'mts')
    OR (t.slug IN ('zero-trust','soc-automation','ddos-protection','post-quantum-crypto') AND o.slug = 'palo-alto')
    OR (t.slug IN ('observability','rag-platform','recommendation-ai','churn-prediction','fraud-ml') AND o.slug = 'snowflake')
ON CONFLICT DO NOTHING;

COMMIT;
