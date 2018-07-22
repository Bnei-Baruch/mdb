-- meaningful terms histogram
select
  tr.term_taxonomy_id,
  t.term_id,
  t.slug,
  count(distinct tr.object_id)
from wp_term_relationships tr
  inner join wp_posts p on tr.object_id = p.ID and p.post_type = 'post'
  inner join wp_term_taxonomy tt on tr.term_taxonomy_id = tt.term_taxonomy_id
  inner join wp_terms t on tt.term_id = t.term_id
group by tr.term_taxonomy_id
having count(distinct tr.object_id) > 100
order by count(distinct tr.object_id);

-- slug base url
select concat('https://laitman.ru/', t.slug, '/', p.ID, '.html')
from wp_posts p
  inner join wp_term_relationships tr on p.ID = tr.object_id
  inner join wp_term_taxonomy tt on tr.term_taxonomy_id = tt.term_taxonomy_id
  inner join wp_terms t on tt.term_id = t.term_id
where p.post_type = 'post' and p.ID = 198975;

-- term samples
select p.ID, concat('https://laitman.ru/', t.slug, '/', p.ID, '.html')
from wp_posts p
  inner join wp_term_relationships tr on p.ID = tr.object_id
  inner join wp_term_taxonomy tt on tr.term_taxonomy_id = tt.term_taxonomy_id
  inner join wp_terms t on tt.term_id = t.term_id
where p.post_type = 'post' and t.slug = 'group'
limit 10;

select p.ID, concat('https://laitman.ru/', t.slug, '/', p.ID, '.html')
from wp_posts p
  inner join wp_term_relationships tr on p.ID = tr.object_id
  inner join wp_term_taxonomy tt on tr.term_taxonomy_id = tt.term_taxonomy_id
  inner join wp_terms t on tt.term_id = t.term_id
where p.post_type = 'post' and t.slug = 'ezhednevny-urok'
  order by length(p.post_content) desc
limit 100;

select
  p.ID,
  p.post_content,
  p.post_title,
  group_concat(distinct t.slug) as terms
from wp_posts p
  inner join wp_term_relationships tr on p.ID = tr.object_id
  inner join wp_term_taxonomy tt on tr.term_taxonomy_id = tt.term_taxonomy_id
  inner join wp_terms t on tt.term_id = t.term_id
where p.post_type = 'post'
group by p.ID
limit 100;

select
  p.ID,
  p.post_type,
  p.post_title
from wp_postmeta m
  inner join wp_posts p on m.post_id = p.ID and m.meta_key = 'podPressMedia'
order by p.post_title, p.post_date;
