// Golf Gamez PWA JavaScript - Enhancement and PWA features

class GolfGamezPWA {
  constructor() {
    this.isOnline = navigator.onLine;
    this.installPrompt = null;
    this.swRegistration = null;

    this.init();
  }

  async init() {
    console.log('Golf Gamez PWA: Initializing...');

    // Register service worker
    await this.registerServiceWorker();

    // Set up PWA features
    this.setupNetworkDetection();
    this.setupInstallPrompt();
    this.setupNotifications();
    this.setupBackgroundSync();

    // Add PWA UI enhancements
    this.addPWAUI();

    console.log('Golf Gamez PWA: Initialization complete');
  }

  // Service Worker Registration
  async registerServiceWorker() {
    if ('serviceWorker' in navigator) {
      try {
        console.log('Golf Gamez PWA WEB: Registering service worker...');

        this.swRegistration = await navigator.serviceWorker.register('/web/sw.js', {
          scope: '/web/'
        });

        console.log('Golf Gamez PWA: Service worker registered successfully', this.swRegistration);

        // Handle service worker updates
        this.swRegistration.addEventListener('updatefound', () => {
          const newWorker = this.swRegistration.installing;

          newWorker.addEventListener('statechange', () => {
            if (newWorker.state === 'installed') {
              if (navigator.serviceWorker.controller) {
                // New service worker available
                this.showUpdateAvailable();
              }
            }
          });
        });

      } catch (error) {
        console.error('Golf Gamez PWA: Service worker registration failed', error);
      }
    }
  }

  // Network Detection
  setupNetworkDetection() {
    // Initial state
    this.updateNetworkStatus();

    // Listen for network changes
    window.addEventListener('online', () => {
      this.isOnline = true;
      this.updateNetworkStatus();
      this.syncWhenOnline();
    });

    window.addEventListener('offline', () => {
      this.isOnline = false;
      this.updateNetworkStatus();
    });
  }

  updateNetworkStatus() {
    // Add network status indicator to the page
    let statusIndicator = document.getElementById('network-status');

    if (!statusIndicator) {
      statusIndicator = document.createElement('div');
      statusIndicator.id = 'network-status';
      statusIndicator.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 1000;
        padding: 8px;
        text-align: center;
        font-size: 14px;
        font-weight: 600;
        transform: translateY(-100%);
        transition: transform 0.3s ease;
      `;
      document.body.appendChild(statusIndicator);
    }

    if (!this.isOnline) {
      statusIndicator.textContent = '⚠️ You\'re offline - some features may be limited';
      statusIndicator.style.backgroundColor = '#ff9800';
      statusIndicator.style.color = '#fff';
      statusIndicator.style.transform = 'translateY(0)';
    } else {
      statusIndicator.style.transform = 'translateY(-100%)';

      // Show brief "back online" message
      setTimeout(() => {
        if (this.isOnline) {
          statusIndicator.textContent = '✅ Back online';
          statusIndicator.style.backgroundColor = '#4caf50';
          statusIndicator.style.color = '#fff';
          statusIndicator.style.transform = 'translateY(0)';

          setTimeout(() => {
            statusIndicator.style.transform = 'translateY(-100%)';
          }, 3000);
        }
      }, 100);
    }
  }

  // Install Prompt Handling
  setupInstallPrompt() {
    window.addEventListener('beforeinstallprompt', (event) => {
      console.log('Golf Gamez PWA: Install prompt available');

      // Prevent the mini-infobar from appearing
      event.preventDefault();

      // Store the event for later use
      this.installPrompt = event;

      // Show custom install button
      this.showInstallButton();
    });

    // Handle successful installation
    window.addEventListener('appinstalled', () => {
      console.log('Golf Gamez PWA: App installed successfully');
      this.installPrompt = null;
      this.hideInstallButton();

      // Show success message
      this.showMessage('Golf Gamez installed successfully! 🎉', 'success');
    });
  }

  showInstallButton() {
    let installButton = document.getElementById('pwa-install-button');

    if (!installButton) {
      installButton = document.createElement('button');
      installButton.id = 'pwa-install-button';
      installButton.textContent = '📱 Install Golf Gamez';
      installButton.style.cssText = `
        position: fixed;
        bottom: 20px;
        right: 20px;
        background: #2e7d32;
        color: white;
        border: none;
        padding: 12px 20px;
        border-radius: 25px;
        font-weight: 600;
        cursor: pointer;
        box-shadow: 0 4px 12px rgba(46, 125, 50, 0.3);
        z-index: 1000;
        transition: all 0.3s ease;
        transform: scale(0);
      `;

      installButton.addEventListener('click', () => this.promptInstall());
      document.body.appendChild(installButton);
    }

    // Animate in
    setTimeout(() => {
      installButton.style.transform = 'scale(1)';
    }, 100);
  }

  hideInstallButton() {
    const installButton = document.getElementById('pwa-install-button');
    if (installButton) {
      installButton.style.transform = 'scale(0)';
      setTimeout(() => {
        installButton.remove();
      }, 300);
    }
  }

  async promptInstall() {
    if (!this.installPrompt) return;

    try {
      // Show the install prompt
      const result = await this.installPrompt.prompt();
      console.log('Golf Gamez PWA: Install prompt result', result);

      // Reset the install prompt
      this.installPrompt = null;
      this.hideInstallButton();

    } catch (error) {
      console.error('Golf Gamez PWA: Install prompt failed', error);
    }
  }

  // Notifications
  async setupNotifications() {
    if ('Notification' in window) {
      const permission = await Notification.requestPermission();
      console.log('Golf Gamez PWA: Notification permission', permission);
    }
  }

  showNotification(title, options = {}) {
    if ('Notification' in window && Notification.permission === 'granted') {
      const notification = new Notification(title, {
        icon: '/web/static/icon-192.png',
        badge: '/web/static/icon-72.png',
        ...options
      });

      // Auto-close after 5 seconds
      setTimeout(() => notification.close(), 5000);

      return notification;
    }
  }

  // Background Sync
  setupBackgroundSync() {
    if (this.swRegistration && 'sync' in window.ServiceWorkerRegistration.prototype) {
      console.log('Golf Gamez PWA: Background sync available');

      // Listen for sync events from service worker
      navigator.serviceWorker.addEventListener('message', (event) => {
        const { type, data } = event.data;

        switch (type) {
          case 'SCORE_SYNCED':
            this.handleScoreSynced(data);
            break;
        }
      });
    }
  }

  async syncWhenOnline() {
    if (this.swRegistration && 'sync' in window.ServiceWorkerRegistration.prototype) {
      try {
        await this.swRegistration.sync.register('score-sync');
        console.log('Golf Gamez PWA: Background sync registered');
      } catch (error) {
        console.error('Golf Gamez PWA: Background sync registration failed', error);
      }
    }
  }

  handleScoreSynced(data) {
    if (data.success) {
      this.showMessage('Scores synced successfully! ✅', 'success');
    } else {
      this.showMessage('Failed to sync some scores ❌', 'error');
    }
  }

  // PWA UI Enhancements
  addPWAUI() {
    // Add viewport height fix for mobile browsers
    this.addViewportHeightFix();

    // Add touch feedback
    this.addTouchFeedback();

    // Add keyboard navigation enhancements
    this.addKeyboardNavigation();

    // Add pull-to-refresh
    this.addPullToRefresh();
  }

  addViewportHeightFix() {
    // Fix for mobile browsers that change viewport height
    const setVH = () => {
      const vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty('--vh', `${vh}px`);
    };

    setVH();
    window.addEventListener('resize', setVH);
    window.addEventListener('orientationchange', setVH);
  }

  addTouchFeedback() {
    // Add subtle touch feedback to interactive elements
    const style = document.createElement('style');
    style.textContent = `
      .touch-feedback {
        transition: transform 0.1s ease, box-shadow 0.1s ease;
      }
      .touch-feedback:active {
        transform: scale(0.98);
        box-shadow: inset 0 2px 4px rgba(0,0,0,0.2);
      }
    `;
    document.head.appendChild(style);

    // Apply to buttons and interactive elements
    document.addEventListener('DOMContentLoaded', () => {
      const interactiveElements = document.querySelectorAll('button, .primary-button, .secondary-button, input[type="submit"]');
      interactiveElements.forEach(el => el.classList.add('touch-feedback'));
    });
  }

  addKeyboardNavigation() {
    // Improve keyboard navigation
    document.addEventListener('keydown', (event) => {
      // Skip navigation for input elements
      if (event.target.matches('input, textarea, select')) return;

      switch (event.key) {
        case 'Enter':
        case ' ':
          // Activate focused button
          if (event.target.matches('button')) {
            event.target.click();
            event.preventDefault();
          }
          break;
      }
    });
  }

  addPullToRefresh() {
    if (!('ontouchstart' in window)) return; // Desktop browsers

    let startY = 0;
    let currentY = 0;
    let pulling = false;
    let pullThreshold = 100;

    const pullToRefreshElement = document.createElement('div');
    pullToRefreshElement.id = 'pull-to-refresh';
    pullToRefreshElement.style.cssText = `
      position: fixed;
      top: -60px;
      left: 0;
      right: 0;
      height: 60px;
      background: linear-gradient(to bottom, #2e7d32, #4caf50);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      z-index: 999;
      transition: transform 0.3s ease;
    `;
    pullToRefreshElement.textContent = '⬇️ Pull to refresh';
    document.body.appendChild(pullToRefreshElement);

    document.addEventListener('touchstart', (e) => {
      if (window.scrollY === 0) {
        startY = e.touches[0].clientY;
        pulling = true;
      }
    });

    document.addEventListener('touchmove', (e) => {
      if (!pulling) return;

      currentY = e.touches[0].clientY;
      const pullDistance = currentY - startY;

      if (pullDistance > 0) {
        e.preventDefault();
        const pullProgress = Math.min(pullDistance / pullThreshold, 1);

        pullToRefreshElement.style.transform = `translateY(${pullDistance * 0.5}px)`;

        if (pullProgress >= 1) {
          pullToRefreshElement.textContent = '🔄 Release to refresh';
        } else {
          pullToRefreshElement.textContent = '⬇️ Pull to refresh';
        }
      }
    });

    document.addEventListener('touchend', () => {
      if (!pulling) return;

      const pullDistance = currentY - startY;
      pulling = false;

      if (pullDistance >= pullThreshold) {
        // Trigger refresh
        pullToRefreshElement.textContent = '🔄 Refreshing...';
        setTimeout(() => {
          window.location.reload();
        }, 500);
      } else {
        // Reset
        pullToRefreshElement.style.transform = 'translateY(-60px)';
        pullToRefreshElement.textContent = '⬇️ Pull to refresh';
      }
    });
  }

  // Utility Methods
  showMessage(message, type = 'info') {
    const messageEl = document.createElement('div');
    messageEl.style.cssText = `
      position: fixed;
      top: 20px;
      left: 50%;
      transform: translateX(-50%) translateY(-100px);
      background: ${type === 'success' ? '#4caf50' : type === 'error' ? '#f44336' : '#2196f3'};
      color: white;
      padding: 12px 24px;
      border-radius: 25px;
      font-weight: 600;
      z-index: 1001;
      transition: transform 0.3s ease;
      box-shadow: 0 4px 12px rgba(0,0,0,0.3);
    `;
    messageEl.textContent = message;
    document.body.appendChild(messageEl);

    // Animate in
    setTimeout(() => {
      messageEl.style.transform = 'translateX(-50%) translateY(0)';
    }, 100);

    // Remove after 4 seconds
    setTimeout(() => {
      messageEl.style.transform = 'translateX(-50%) translateY(-100px)';
      setTimeout(() => messageEl.remove(), 300);
    }, 4000);
  }

  showUpdateAvailable() {
    const updateBanner = document.createElement('div');
    updateBanner.style.cssText = `
      position: fixed;
      bottom: 0;
      left: 0;
      right: 0;
      background: #ff9800;
      color: white;
      padding: 16px;
      text-align: center;
      z-index: 1000;
      transform: translateY(100%);
      transition: transform 0.3s ease;
    `;
    updateBanner.innerHTML = `
      <p style="margin: 0 0 8px 0; font-weight: 600;">New version available!</p>
      <button id="update-app" style="
        background: white;
        color: #ff9800;
        border: none;
        padding: 8px 16px;
        border-radius: 4px;
        font-weight: 600;
        cursor: pointer;
        margin-right: 8px;
      ">Update Now</button>
      <button id="dismiss-update" style="
        background: transparent;
        color: white;
        border: 1px solid white;
        padding: 8px 16px;
        border-radius: 4px;
        font-weight: 600;
        cursor: pointer;
      ">Later</button>
    `;

    document.body.appendChild(updateBanner);

    // Animate in
    setTimeout(() => {
      updateBanner.style.transform = 'translateY(0)';
    }, 100);

    // Handle update
    document.getElementById('update-app').addEventListener('click', () => {
      this.swRegistration.waiting.postMessage({ type: 'SKIP_WAITING' });
      window.location.reload();
    });

    // Handle dismiss
    document.getElementById('dismiss-update').addEventListener('click', () => {
      updateBanner.style.transform = 'translateY(100%)';
      setTimeout(() => updateBanner.remove(), 300);
    });
  }

  // API Enhancement Methods
  async enhancedFetch(url, options = {}) {
    try {
      const response = await fetch(url, options);

      // Check if response indicates offline mode
      if (response.headers.get('X-Golf-Gamez-Offline') === 'true') {
        this.showMessage('Using cached data (offline mode)', 'warning');
      }

      return response;
    } catch (error) {
      // If we're offline and this is a score submission, cache it
      if (!this.isOnline && options.method === 'POST' && url.includes('/scores')) {
        await this.cacheScoreForSync(url, options);
        this.showMessage('Score saved - will sync when online', 'info');

        // Return a mock success response
        return new Response(JSON.stringify({ cached: true }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        });
      }

      throw error;
    }
  }

  async cacheScoreForSync(url, options) {
    if (this.swRegistration) {
      navigator.serviceWorker.controller.postMessage({
        type: 'CACHE_SCORE',
        data: {
          url,
          method: options.method,
          headers: options.headers,
          body: options.body,
          timestamp: Date.now()
        }
      });
    }
  }
}

// Initialize PWA features when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    window.golfGamezPWA = new GolfGamezPWA();
  });
} else {
  window.golfGamezPWA = new GolfGamezPWA();
}

// Export for use in Go WebAssembly code
window.GolfGamezPWA = GolfGamezPWA;