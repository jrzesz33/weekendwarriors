// Golf Gamez Service Worker
// Provides offline functionality and caching for the PWA

const CACHE_NAME = 'golf-gamez-v1.0.0';
const API_CACHE_NAME = 'golf-gamez-api-v1.0.0';

// Static assets to cache
const STATIC_ASSETS = [
  '/',
  '/web/css/app.css',
  '/web/js/app.js',
  '/web/static/icon-192.png',
  '/web/static/icon-512.png',
  '/web/static/icon-180.png',
  '/app.wasm',
  '/wasm_exec.js',
  '/manifest.json'
];

// API endpoints that can be cached
const CACHEABLE_API_PATTERNS = [
  /\/v1\/games\/[^\/]+$/,           // Game info
  /\/v1\/games\/[^\/]+\/players$/,  // Player list
  /\/v1\/spectate\/[^\/]+$/         // Spectator views
];

// Install event - cache static assets
self.addEventListener('install', event => {
  console.log('Service Worker: Installing...');

  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => {
        console.log('Service Worker: Caching static assets');
        return cache.addAll(STATIC_ASSETS);
      })
      .then(() => {
        console.log('Service Worker: Installation complete');
        // Force activation of new service worker
        return self.skipWaiting();
      })
      .catch(error => {
        console.error('Service Worker: Installation failed', error);
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', event => {
  console.log('Service Worker: Activating...');

  event.waitUntil(
    caches.keys()
      .then(cacheNames => {
        return Promise.all(
          cacheNames.map(cacheName => {
            if (cacheName !== CACHE_NAME && cacheName !== API_CACHE_NAME) {
              console.log('Service Worker: Deleting old cache', cacheName);
              return caches.delete(cacheName);
            }
          })
        );
      })
      .then(() => {
        console.log('Service Worker: Activation complete');
        // Take control of all clients immediately
        return self.clients.claim();
      })
  );
});

// Fetch event - implement caching strategies
self.addEventListener('fetch', event => {
  const request = event.request;
  const url = new URL(request.url);

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return;
  }

  // Skip chrome-extension and browser internal requests
  if (url.protocol !== 'http:' && url.protocol !== 'https:') {
    return;
  }

  // Handle different types of requests
  if (isStaticAsset(request)) {
    event.respondWith(handleStaticAsset(request));
  } else if (isAPIRequest(request)) {
    event.respondWith(handleAPIRequest(request));
  } else if (isNavigationRequest(request)) {
    event.respondWith(handleNavigationRequest(request));
  }
});

// Check if request is for a static asset
function isStaticAsset(request) {
  const url = new URL(request.url);
  return STATIC_ASSETS.some(asset => url.pathname.endsWith(asset)) ||
         url.pathname.includes('/web/') ||
         url.pathname.endsWith('.wasm') ||
         url.pathname.endsWith('.js') ||
         url.pathname.endsWith('.css') ||
         url.pathname.endsWith('.png') ||
         url.pathname.endsWith('.json');
}

// Check if request is for API
function isAPIRequest(request) {
  const url = new URL(request.url);
  return url.pathname.startsWith('/v1/') || url.pathname.startsWith('/api/');
}

// Check if request is a navigation request
function isNavigationRequest(request) {
  return request.mode === 'navigate' ||
         (request.method === 'GET' && request.headers.get('accept').includes('text/html'));
}

// Handle static assets with cache-first strategy
async function handleStaticAsset(request) {
  try {
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }

    const networkResponse = await fetch(request);
    if (networkResponse.ok) {
      const cache = await caches.open(CACHE_NAME);
      cache.put(request, networkResponse.clone());
    }
    return networkResponse;
  } catch (error) {
    console.error('Service Worker: Failed to fetch static asset', error);

    // Return offline fallback for critical assets
    if (request.url.includes('app.css')) {
      return new Response('/* Offline mode - styles unavailable */', {
        headers: { 'Content-Type': 'text/css' }
      });
    }

    throw error;
  }
}

// Handle API requests with network-first strategy
async function handleAPIRequest(request) {
  const url = new URL(request.url);

  try {
    // Try network first
    const networkResponse = await fetch(request);

    if (networkResponse.ok) {
      // Cache GET requests for specific endpoints
      if (request.method === 'GET' && isCacheableAPI(url)) {
        const cache = await caches.open(API_CACHE_NAME);
        cache.put(request, networkResponse.clone());
      }
      return networkResponse;
    }

    throw new Error(`API request failed: ${networkResponse.status}`);
  } catch (error) {
    console.warn('Service Worker: Network request failed, trying cache', error);

    // Fall back to cache for GET requests
    if (request.method === 'GET') {
      const cachedResponse = await caches.match(request);
      if (cachedResponse) {
        // Add offline indicator header
        const response = new Response(cachedResponse.body, {
          status: cachedResponse.status,
          statusText: cachedResponse.statusText,
          headers: {
            ...Object.fromEntries(cachedResponse.headers.entries()),
            'X-Golf-Gamez-Offline': 'true'
          }
        });
        return response;
      }
    }

    // Return offline error response
    return new Response(
      JSON.stringify({
        error: {
          code: 'offline_error',
          message: 'Unable to connect to server. Please check your internet connection.',
          offline: true
        }
      }),
      {
        status: 503,
        statusText: 'Service Unavailable',
        headers: {
          'Content-Type': 'application/json',
          'X-Golf-Gamez-Offline': 'true'
        }
      }
    );
  }
}

// Handle navigation requests
async function handleNavigationRequest(request) {
  try {
    return await fetch(request);
  } catch (error) {
    console.warn('Service Worker: Navigation request failed, serving app shell', error);

    // Serve the app shell for offline navigation
    const cache = await caches.open(CACHE_NAME);
    const appShell = await cache.match('/');

    if (appShell) {
      return appShell;
    }

    // Fallback offline page
    return new Response(`
      <!DOCTYPE html>
      <html>
      <head>
        <title>Golf Gamez - Offline</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <style>
          body {
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            text-align: center;
            padding: 2rem;
            background: #f5f5f5;
          }
          .offline-container {
            max-width: 400px;
            margin: 0 auto;
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
          }
          h1 { color: #2e7d32; }
          .icon { font-size: 3rem; margin-bottom: 1rem; }
          button {
            background: #2e7d32;
            color: white;
            border: none;
            padding: 1rem 2rem;
            border-radius: 4px;
            font-size: 1rem;
            cursor: pointer;
            margin-top: 1rem;
          }
          button:hover { background: #1b5e20; }
        </style>
      </head>
      <body>
        <div class="offline-container">
          <div class="icon">⛳</div>
          <h1>Golf Gamez</h1>
          <h2>You're Offline</h2>
          <p>Please check your internet connection and try again.</p>
          <button onclick="location.reload()">Try Again</button>
        </div>
      </body>
      </html>
    `, {
      headers: { 'Content-Type': 'text/html' }
    });
  }
}

// Check if API endpoint should be cached
function isCacheableAPI(url) {
  return CACHEABLE_API_PATTERNS.some(pattern => pattern.test(url.pathname));
}

// Background sync for score updates when back online
self.addEventListener('sync', event => {
  console.log('Service Worker: Background sync event', event.tag);

  if (event.tag === 'score-sync') {
    event.waitUntil(syncPendingScores());
  }
});

// Sync pending scores when connection is restored
async function syncPendingScores() {
  try {
    // Get pending scores from IndexedDB
    const pendingScores = await getPendingScores();

    for (const score of pendingScores) {
      try {
        const response = await fetch(score.url, {
          method: 'POST',
          headers: score.headers,
          body: score.body
        });

        if (response.ok) {
          // Remove from pending scores
          await removePendingScore(score.id);

          // Notify all clients of successful sync
          const clients = await self.clients.matchAll();
          clients.forEach(client => {
            client.postMessage({
              type: 'SCORE_SYNCED',
              data: { scoreId: score.id, success: true }
            });
          });
        }
      } catch (error) {
        console.error('Service Worker: Failed to sync score', error);
      }
    }
  } catch (error) {
    console.error('Service Worker: Background sync failed', error);
  }
}

// Placeholder functions for IndexedDB operations
// These would be implemented with proper IndexedDB handling
async function getPendingScores() {
  // Return pending scores from IndexedDB
  return [];
}

async function removePendingScore(scoreId) {
  // Remove score from IndexedDB
  return true;
}

// Handle messages from the main application
self.addEventListener('message', event => {
  const { type, data } = event.data;

  switch (type) {
    case 'SKIP_WAITING':
      self.skipWaiting();
      break;

    case 'CACHE_SCORE':
      // Store score for later sync
      event.waitUntil(cacheScoreForSync(data));
      break;

    case 'CLEAR_CACHE':
      event.waitUntil(clearAllCaches());
      break;

    default:
      console.log('Service Worker: Unknown message type', type);
  }
});

// Cache score for background sync
async function cacheScoreForSync(scoreData) {
  try {
    // Store in IndexedDB for background sync
    console.log('Service Worker: Caching score for sync', scoreData);

    // Register for background sync
    await self.registration.sync.register('score-sync');
  } catch (error) {
    console.error('Service Worker: Failed to cache score for sync', error);
  }
}

// Clear all caches
async function clearAllCaches() {
  try {
    const cacheNames = await caches.keys();
    await Promise.all(cacheNames.map(cacheName => caches.delete(cacheName)));
    console.log('Service Worker: All caches cleared');
  } catch (error) {
    console.error('Service Worker: Failed to clear caches', error);
  }
}

// Push notification handling (for future use)
self.addEventListener('push', event => {
  if (!event.data) return;

  const data = event.data.json();
  const options = {
    body: data.body,
    icon: '/web/static/icon-192.png',
    badge: '/web/static/icon-72.png',
    tag: data.tag || 'golf-gamez-notification',
    data: data.data,
    actions: data.actions || [],
    requireInteraction: data.requireInteraction || false
  };

  event.waitUntil(
    self.registration.showNotification(data.title, options)
  );
});

// Handle notification clicks
self.addEventListener('notificationclick', event => {
  event.notification.close();

  const urlToOpen = event.notification.data?.url || '/';

  event.waitUntil(
    self.clients.matchAll({ type: 'window' })
      .then(clients => {
        // Check if there's already a window open
        for (const client of clients) {
          if (client.url === urlToOpen && 'focus' in client) {
            return client.focus();
          }
        }

        // Open a new window
        if (self.clients.openWindow) {
          return self.clients.openWindow(urlToOpen);
        }
      })
  );
});

console.log('Service Worker: Script loaded');