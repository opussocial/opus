class HttpClient {
  constructor(baseURL = '') {
    this.baseURL = baseURL;
    this.authToken = localStorage.getItem("auth:token");
  }

  get(endpoint, params = {}) {
    const url = this.buildUrl(endpoint, params);
    return this.request('GET', url);
  }

  post(endpoint, body) {
    return this.request('POST', this.buildUrl(endpoint), body);
  }

  put(endpoint, body) {
    return this.request('PUT', this.buildUrl(endpoint), body);
  }

  delete(endpoint) {
    return this.request('DELETE', this.buildUrl(endpoint));
  }

  buildUrl(endpoint, params = {}) {
    const url = new URL(this.baseURL + endpoint);
    Object.entries(params).forEach(([key, value]) => {
      url.searchParams.append(key, value);
    });
    return url.toString();
  }

  request(method, url, body = null) {
    const options = {
      method,
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': 'Bearer: ' + this.authToken 
      },
      ...(body && { body: JSON.stringify(body) }),
    };

    return rxjs.from(fetch(url, options)).pipe(
      rxjs.operators.switchMap(response => {
        if (!response.ok) {
          return rxjs.throwError(() => new Error(`HTTP ${response.status}`));
        }
        return rxjs.from(response.json());
      }),
      rxjs.operators.catchError(err => {
        console.error('HTTP Error:', err);
        return rxjs.throwError(() => err);
      })
    );
  }
}

// RxJS Event Bus
class EventBus {
  constructor(hooks) {
    this.subjects = new Map(); // Stores RxJS Subjects for each event
    this.hooks = hooks;
  }

  // Emit an event
  publish(event, data) {
    console.log("fired:", event);
    if (this.hooks.onBefore) {
      this.hooks.onBefore(event, data);
    }
    if (!this.subjects.has(event)) {
      this.subjects.set(event, new rxjs.Subject());
    }
    this.subjects.get(event).next(data || {});
    if (this.hooks.onAfter) {
      this.hooks.onAfter(event, data);
    }
  }

  // Subscribe to an event
  subscribe(eventName, handler) {
    if (!this.subjects.has(eventName)) {
      this.subjects.set(eventName, new rxjs.Subject());
    }
    if (this.hooks.onSubscribe) {
      this.hooks.onSubscribe(eventName);
    }
    return this.subjects.get(eventName).subscribe(handler);
  }

  // Unsubscribe (cleanup)
  cancel(subscription) {
    if (subscription) subscription.unsubscribe();
  }

  // Clear all subscriptions for an event
  clear(eventName) {
    if (this.subjects.has(eventName)) {
      this.subjects.get(eventName).complete();
      this.subjects.delete(eventName);
    }
  }
}

class Component {
  constructor(behavior, store) {
    if (behavior) {
      this.behavior = behavior;
    }
    if (store) {
      this.store = PetiteVue.reactive(store);
    }
    if (router) {
      this.router = router;
    }
  }

  subscribe(bus) {
    this.subscriptions = this.behavior.subscriptions.map((i) => { bus.subscribe });
    return this;
  }

  route() {
    const router = Router(this.behavior.routes)
    router.init();
    return this;
  }

  mount(elementId) {
    const scope = {
      store: this.store,
      init: this.behavior.init,
      fire(event, data) {
        bus.publish(event, data);
      },
      screen(screen) {
        return useScreen(screen);
      },
    };
    PetiteVue.createApp(this.scope).mount(elementId)
    return this;
  }
}

// Read query string parameters
function getQueryParam(name) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(name);
}

// Get all query parameters as object
function getAllQueryParams() {
  const params = {};
  const urlParams = new URLSearchParams(window.location.search);
  
  for (const [key, value] of urlParams.entries()) {
    params[key] = value;
  }
  
  return params;
}

// Example usage:
// const id = getQueryParam('id');
// const allParams = getAllQueryParams();

// Read cookie
function getCookie(name) {
  const nameEQ = name + "=";
  const cookies = document.cookie.split(';');
  
  for (let i = 0; i < cookies.length; i++) {
    let cookie = cookies[i].trim();
    if (cookie.indexOf(nameEQ) === 0) {
      return decodeURIComponent(cookie.substring(nameEQ.length));
    }
  }
  return null;
}

// Write cookie
function setCookie(name, value, days = 30, path = '/') {
  let expires = '';
  if (days) {
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    expires = '; expires=' + date.toUTCString();
  }
  document.cookie = name + '=' + encodeURIComponent(value) + expires + '; path=' + path;
}

// Delete cookie
function deleteCookie(name, path = '/') {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=' + path;
}

// Example usage:
// setCookie('theme', 'dark', 7);
// const theme = getCookie('theme');
// deleteCookie('theme');
