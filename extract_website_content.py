import requests
from bs4 import BeautifulSoup
from bs4.element import Comment
from urllib.parse import urljoin, urlparse, urlunparse
import sys
import xml.etree.ElementTree as ET
import time
import json
import uuid
import logging
from collections import Counter
import re

# Replace with your sitemap URL
sitemap_url = 'https://stevenlawton.com/sitemap.xml'

# Web service endpoint to post JSON data
web_service_url = 'http://localhost:8080/api/receive_data'

# Generate a unique extract ID
extract_id = str(uuid.uuid4())

visited = set()
to_visit = []

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Function to parse sitemap and extract URLs
def extract_urls_from_sitemap(sitemap_url):
    try:
        response = requests.get(sitemap_url, timeout=10)
        if response.status_code != 200:
            logging.error(f'Failed to fetch sitemap: {response.status_code}')
            return []
        sitemap_content = response.content
        root = ET.fromstring(sitemap_content)

        urls = []
        # If it's a sitemap index, get all nested sitemap URLs
        for elem in root.iter():
            if elem.tag.endswith("sitemapindex"):
                for sitemap in root.findall(".//{*}sitemap/{*}loc"):
                    sitemap_url = sitemap.text.strip()
                    urls.extend(extract_urls_from_sitemap(sitemap_url))
            elif elem.tag.endswith("loc"):  # Extract all <loc> elements from regular sitemaps
                urls.append(elem.text.strip())

        return urls
    except requests.exceptions.RequestException as e:
        logging.error(f'Failed to fetch sitemap: {e}')
        return []

# Function to check if a tag is visible
def tag_visible(element):
    if element.parent.name in ['style', 'script', 'head', 'title', 'meta', '[document]']:
        return False
    if isinstance(element, Comment):
        return False
    return True

# Function to determine if a link is an HTML page
def is_html_page(url):
    excluded_extensions = ('.js', '.css', '.jpg', '.jpeg', '.png', '.gif', '.pdf', '.zip', '.rar', '.mp4', '.mp3')
    return not url.lower().endswith(excluded_extensions)

# Function to determine if a response contains an HTML page
def is_html_response(response):
    content_type = response.headers.get('Content-Type', '').lower()
    return 'text/html' in content_type

def normalise_url(url):
    parsed_url = urlparse(url)
    # Remove the fragment part of the URL
    normalised_url = urlunparse(parsed_url._replace(fragment=''))
    return normalised_url


# Extract URLs from the sitemap
logging.info("Extracting URLs from sitemap...")
to_visit = extract_urls_from_sitemap(sitemap_url)
total_urls = len(to_visit)
logging.info(f"Found {total_urls} URLs to process.\n")

# Define the domain for comparison
domain = urlparse(sitemap_url).netloc

# Process each URL
current_index = 0
while to_visit:
    url = to_visit.pop(0)
    if url in visited:
        continue
    visited.add(url)

    # Skip URLs that are non-HTML files
    if not is_html_page(url):
        logging.info(f"Skipping non-HTML URL: {url}")
        continue

    # Measure page load time
    logging.info(f"Processing URL {current_index + 1}/{total_urls}: {url}")
    current_index += 1
    start_time = time.time()
    try:
        response = requests.get(url, timeout=10)
    except requests.exceptions.RequestException as e:
        logging.error(f'Failed to fetch {url}: {e}')
        continue
    load_time = time.time() - start_time

    # Skip non-HTML responses
    if not is_html_response(response):
        logging.info(f"Skipping non-HTML response from URL: {url}")
        continue

    if response.status_code != 200:
        logging.warning(f'Non-200 response for {url}: {response.status_code}')
        continue

    soup = BeautifulSoup(response.content, 'html.parser')

    # Ensure soup.html is not None
    if soup.html is None:
        logging.warning(f"Skipping URL due to missing HTML content: {url}")
        continue

    # Extract internal and external links, and add new internal links to 'to_visit'
    for link in soup.find_all('a', href=True):
        href = link['href']
        full_url = normalise_url(urljoin(url, href))
        link_netloc = urlparse(full_url).netloc

        if link_netloc == domain and is_html_page(full_url):
            if full_url not in visited and full_url not in to_visit:
                to_visit.append(full_url)

    # Extract page language
    lang = soup.html.get('lang', 'No Language Specified')

    # Extract page title
    title = soup.title.string if soup.title else 'No Title'

    # Extract meta tags
    meta_tags = {}
    for meta in soup.find_all('meta'):
        name = meta.get('name', '').lower()
        prop = meta.get('property', '').lower()
        if name or prop:
            key = name if name else prop
            meta_tags[key] = meta.get('content', '')

    # Extract canonical tag
    canonical = soup.find('link', rel='canonical')
    canonical_url = canonical['href'] if canonical else 'No Canonical URL'
    is_canonical_correct = (canonical_url == url)

    # Extract H tags
    h_tags = {}
    for i in range(1, 7):
        h_tag = f'h{i}'
        h_tags[h_tag] = [h.get_text(strip=True) for h in soup.find_all(h_tag)]

    # Extract visible text from the page
    texts = soup.findAll(string=True)
    visible_texts = filter(tag_visible, texts)
    body = u" ".join(t.strip() for t in visible_texts)

    # Calculate word count
    word_count = len(body.split())

    # Extract images and their alt attributes
    images = [{'src': img.get('src'), 'alt': img.get('alt', 'No alt attribute'),
               'width': img.get('width'), 'height': img.get('height')} for img in soup.find_all('img')]

    # Extract internal and external links, and add new internal links to 'to_visit'
    internal_links = []
    internal_links_with_anchor = []
    external_links = []
    domain = urlparse(sitemap_url).netloc

    for link in soup.find_all('a', href=True):
        href = link['href']
        full_url = urljoin(url, href)
        link_netloc = urlparse(full_url).netloc

        if link_netloc == domain and is_html_page(full_url):
            internal_links.append(full_url)
            internal_links_with_anchor.append({'href': full_url, 'anchorText': link.get_text(strip=True)})
            if full_url not in visited and full_url not in to_visit:
                to_visit.append(full_url)
        elif link_netloc != domain:
            external_links.append(full_url)

    # Check for broken links
    broken_links = []
    for link in internal_links + external_links:
        try:
            link_response = requests.head(link, allow_redirects=True, timeout=5)
            if link_response.status_code == 404:
                broken_links.append(link)
        except requests.exceptions.RequestException:
            broken_links.append(link)

    # Extract structured data (JSON-LD)
    structured_data = soup.find_all('script', type='application/ld+json')
    structured_data_content = [sd.get_text(strip=True) for sd in structured_data]
    structured_data_types = []
    for sd in structured_data:
        try:
            data_json = json.loads(sd.get_text(strip=True))
            if isinstance(data_json, dict) and '@type' in data_json:
                structured_data_types.append(data_json['@type'])
        except json.JSONDecodeError:
            continue

    # Extract robots meta tag
    robots_meta = soup.find('meta', attrs={'name': 'robots'})
    robots_content = robots_meta['content'] if robots_meta else 'No Robots Meta Tag'

    # Extract social media tags (Open Graph & Twitter)
    social_tags = {}
    for meta in soup.find_all('meta'):
        if meta.get('property') and (meta.get('property').startswith('og:') or meta.get('property').startswith('twitter:')):
            social_tags[meta.get('property')] = meta.get('content', '')

    # Extract language and locale information
    lang = soup.html.get('lang', 'No Language Specified')
    hreflangs = [link.get('hreflang') for link in soup.find_all('link', rel='alternate') if link.get('hreflang')]

    # Extract breadcrumbs
    breadcrumbs = []
    for item in soup.find_all('nav', attrs={'aria-label': 'breadcrumb'}):
        breadcrumbs.extend([li.get_text(strip=True) for li in item.find_all('li')])

    # Check mobile-friendliness
    viewport_meta = soup.find('meta', attrs={'name': 'viewport'})
    is_mobile_friendly = bool(viewport_meta)

    # Extract external scripts and stylesheets
    external_scripts = [script.get('src') for script in soup.find_all('script', src=True)]
    external_stylesheets = [link.get('href') for link in soup.find_all('link', rel='stylesheet')]

    # Calculate page depth
    page_depth = urlparse(url).path.count('/')

    # Calculate page size
    page_size = len(response.content)

    # Calculate word frequency
    words = re.findall(r'\w+', body.lower())
    word_frequency = Counter(words)
    common_words = word_frequency.most_common(20)

    # Create a dictionary for each page's data
    page_data = {
        'extractId': extract_id,
        'url': url,
        'title': title,
        'titleLength': len(title) if title else 0,
        'metaDescription': meta_tags.get('description', ''),
        'metaDescriptionLength': len(meta_tags.get('description', '')),
        'metaTags': meta_tags,
        'canonicalUrl': canonical_url,
        'isCanonicalCorrect': is_canonical_correct,
        'hTags': h_tags,
        'h1TagCount': len(h_tags.get('h1', [])),
        'wordCount': word_count,
        'pageDepth': page_depth,
        'pageLoadTimeSeconds': round(load_time, 2),
        'pageSizeBytes': page_size,
        'images': images,
        'internalLinks': internal_links,
        'internalLinksWithAnchorText': internal_links_with_anchor,
        'externalLinks': external_links,
        'brokenLinks': broken_links,
        'structuredData': structured_data_content,
        'structuredDataTypes': structured_data_types,
        'robotsMetaTag': robots_content,
        'content': body,
        'commonWords': common_words,
        'socialTags': social_tags,
        'language': lang,
        'hreflangs': hreflangs,
        'breadcrumbs': breadcrumbs,
        'isMobileFriendly': is_mobile_friendly,
        'externalScripts': external_scripts,
        'externalStylesheets': external_stylesheets
    }

    # Post the page data to the web service
    try:
        post_response = requests.post(web_service_url, json=page_data, timeout=10)
        if post_response.status_code == 201:
            logging.info(f"Successfully posted data for {url}\n")
        else:
            logging.error(f'Failed to post data for {url}: {post_response.status_code} - {post_response.text}\n')
    except requests.exceptions.RequestException as e:
        logging.error(f'Failed to post data for {url}: {e}\n')

logging.info("All URLs have been processed.")
