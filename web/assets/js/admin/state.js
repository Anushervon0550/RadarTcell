// ============ ADMIN STATE + ENTITY CONFIG ============

export const AS = {
  me: null,
  entity: 'technologies',
  items: [],
};

export const ENT = {
  technologies: {
    title: 'Технологии', list: '/admin/technologies?limit=200', idKey: 'slug',
    cols: [['index', '№'], ['name', 'Название'], ['slug', 'Slug'], ['trl', 'TRL'], ['trend_slug', 'Тренд']],
    fields: [
      { k: 'slug', label: 'Slug', req: true },
      { k: 'index', label: 'Index', type: 'number' },
      { k: 'name', label: 'Имя', req: true },
      { k: 'trl', label: 'TRL', type: 'number' },
      { k: 'trend_slug', label: 'Trend slug', req: true },
      { k: 'description_short', label: 'Краткое описание', type: 'textarea' },
      { k: 'description_full', label: 'Полное описание', type: 'textarea' },
      { k: 'image_url', label: 'Image URL' },
      { k: 'source_link', label: 'Source link' },
      { k: 'tag_slugs', label: 'Tag slugs (через запятую)', type: 'csv' },
      { k: 'sdg_codes', label: 'SDG codes (через запятую)', type: 'csv' },
      { k: 'organization_slugs', label: 'Org slugs (через запятую)', type: 'csv' },
      { k: 'custom_metric_1', label: 'Custom metric 1', type: 'number' },
      { k: 'custom_metric_2', label: 'Custom metric 2', type: 'number' },
      { k: 'custom_metric_3', label: 'Custom metric 3', type: 'number' },
      { k: 'custom_metric_4', label: 'Custom metric 4', type: 'number' },
    ],
  },
  trends: {
    title: 'Тренды', list: '/admin/trends', idKey: 'slug',
    cols: [['name', 'Имя'], ['slug', 'Slug']],
    fields: [{ k: 'slug', req: true }, { k: 'name', req: true }],
  },
  tags: {
    title: 'Теги', list: '/admin/tags', idKey: 'slug',
    cols: [['title', 'Имя'], ['slug', 'Slug'], ['category', 'Категория']],
    fields: [{ k: 'slug', req: true }, { k: 'title', req: true }, { k: 'category' }],
  },
  organizations: {
    title: 'Организации', list: '/admin/organizations', idKey: 'slug',
    cols: [['name', 'Имя'], ['slug', 'Slug'], ['logo_url', 'Logo']],
    fields: [{ k: 'slug', req: true }, { k: 'name', req: true }, { k: 'logo_url' }],
  },
  metrics: {
    title: 'Метрики', list: '/admin/metrics', idKey: 'id',
    cols: [['name', 'Имя'], ['type', 'Тип'], ['field_key', 'Field']],
    fields: [{ k: 'name', req: true }, { k: 'type', req: true }, { k: 'description' }, { k: 'orderable', type: 'bool' }, { k: 'field_key' }],
  },
  sdgs: {
    title: 'SDG', list: '/admin/sdgs', idKey: 'code',
    cols: [['code', 'Код'], ['title', 'Название']],
    fields: [{ k: 'code', req: true }, { k: 'title', req: true }],
  },
  users: {
    title: 'Пользователи', list: '/admin/users', idKey: 'username',
    cols: [['username', 'Логин'], ['active', 'Активен']],
    fields: [{ k: 'username', req: true }, { k: 'password', req: true, type: 'password' }],
    noEdit: true,
  },
};
