import crypto from 'node:crypto';
import fs from 'node:fs';
import http from 'node:http';

loadDotEnv();

const PORT = Number(process.env.PORT || 3001);
const UPSTREAM_BASE_URL = trimTrailingSlash(
  process.env.UPSTREAM_BASE_URL || 'http://match.biubiubiu.io:8099',
);
const SIGN_KEY = process.env.SIGN_KEY || 'nCRmkbuTucUv';
const SIGN_ALGORITHM = process.env.SIGN_ALGORITHM || 'sha256';
const UPSTREAM_QUERY_MODE = process.env.UPSTREAM_QUERY_MODE || 'json-query';

const KNOWN_ENDPOINTS = new Set([
  '/xdc/get_info',
  '/xdc/get_kill_info',
  '/xd/get_score_info',
  '/xd/get_template',
]);

const server = http.createServer(async (req, res) => {
  try {
    setCorsHeaders(res);

    if (req.method === 'OPTIONS') {
      res.writeHead(204);
      res.end();
      return;
    }

    const requestUrl = new URL(req.url || '/', `http://${req.headers.host || 'localhost'}`);

    if (requestUrl.pathname === '/health') {
      sendJson(res, 200, { status: 'ok' });
      return;
    }

    if (requestUrl.pathname === '/') {
      sendJson(res, 200, {
        name: 'buibuibui-proxy',
        upstream: UPSTREAM_BASE_URL,
        endpoints: [...KNOWN_ENDPOINTS],
        usage: {
          direct: 'GET /xdc/get_info?warId=xxx',
          post: 'POST /xd/get_score_info with JSON body',
        },
      });
      return;
    }

    if (!['GET', 'POST'].includes(req.method || '')) {
      sendJson(res, 405, { error: 'Only GET and POST are supported' });
      return;
    }

    const targetPath = normalizeProxyPath(requestUrl.pathname);
    const query = await readRequestPayload(req, requestUrl);
    const upstreamResponse = await forwardToUpstream(req.method, targetPath, query);

    writeUpstreamResponse(res, upstreamResponse);
  } catch (error) {
    sendJson(res, error.statusCode || 500, {
      error: error.message || 'Internal server error',
    });
  }
});

server.listen(PORT, () => {
  console.log(`Proxy listening on http://localhost:${PORT}`);
  console.log(`Upstream: ${UPSTREAM_BASE_URL}`);
});

function normalizeProxyPath(pathname) {
  const path = pathname.startsWith('/api/') ? pathname.slice(4) : pathname;

  if (!path.startsWith('/')) {
    throw httpError(400, 'Invalid upstream path');
  }

  return path;
}

async function readRequestPayload(req, requestUrl) {
  const queryObject = Object.fromEntries(requestUrl.searchParams.entries());

  if (req.method === 'GET') {
    return coerceObjectValues(queryObject);
  }

  const bodyText = await readBody(req);

  if (!bodyText.trim()) {
    return coerceObjectValues(queryObject);
  }

  try {
    const body = JSON.parse(bodyText);

    if (!body || typeof body !== 'object' || Array.isArray(body)) {
      throw new Error('JSON body must be an object');
    }

    return body;
  } catch (error) {
    throw httpError(400, `Invalid JSON body: ${error.message}`);
  }
}

async function forwardToUpstream(method, path, query) {
  const jsonQuery = stableJsonStringify(query);
  const timestamp = Math.floor(Date.now() / 1000).toString();
  const signature = createSignature(timestamp, path, jsonQuery);
  const headers = {
    accept: 'application/json',
    t: timestamp,
    s: signature,
  };

  const { url, body } = buildUpstreamRequest(method, path, jsonQuery);
  const response = await fetch(url, {
    method,
    headers: body
      ? {
          ...headers,
          'content-type': 'application/json',
        }
      : headers,
    body,
  });

  const contentType = response.headers.get('content-type') || '';
  const responseText = await response.text();

  return {
    status: response.status,
    contentType,
    body: parseJsonIfPossible(responseText),
  };
}

function buildUpstreamRequest(method, path, jsonQuery) {
  const url = new URL(`${UPSTREAM_BASE_URL}${path}`);

  if (method === 'POST') {
    return { url, body: jsonQuery };
  }

  if (jsonQuery !== '{}') {
    if (UPSTREAM_QUERY_MODE === 'json-param') {
      url.searchParams.set('q', jsonQuery);
    } else if (UPSTREAM_QUERY_MODE === 'normal-query') {
      for (const [key, value] of Object.entries(JSON.parse(jsonQuery))) {
        url.searchParams.set(key, String(value));
      }
    } else {
      url.search = jsonQuery;
    }
  }

  return { url, body: undefined };
}

function createSignature(timestamp, path, jsonQuery) {
  const digest = crypto
    .createHash(SIGN_ALGORITHM)
    .update(`${SIGN_KEY}${timestamp}${path}${jsonQuery}`, 'utf8')
    .digest();

  return digest.toString('base64');
}

function stableJsonStringify(value) {
  return JSON.stringify(sortObjectKeys(value));
}

function sortObjectKeys(value) {
  if (Array.isArray(value)) {
    return value.map(sortObjectKeys);
  }

  if (value && typeof value === 'object') {
    return Object.keys(value)
      .sort()
      .reduce((result, key) => {
        result[key] = sortObjectKeys(value[key]);
        return result;
      }, {});
  }

  return value;
}

function coerceObjectValues(object) {
  return Object.fromEntries(
    Object.entries(object).map(([key, value]) => [key, coerceScalar(value)]),
  );
}

function coerceScalar(value) {
  if (value === 'true') return true;
  if (value === 'false') return false;
  if (value !== '' && Number.isFinite(Number(value))) return Number(value);
  return value;
}

function parseJsonIfPossible(text) {
  if (!text) return null;

  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function writeUpstreamResponse(res, upstreamResponse) {
  const isJson = typeof upstreamResponse.body !== 'string';
  const body = isJson
    ? JSON.stringify(upstreamResponse.body)
    : upstreamResponse.body;

  res.writeHead(upstreamResponse.status, {
    'content-type': isJson ? 'application/json; charset=utf-8' : 'text/plain; charset=utf-8',
  });
  res.end(body);
}

function sendJson(res, statusCode, body) {
  res.writeHead(statusCode, {
    'content-type': 'application/json; charset=utf-8',
  });
  res.end(JSON.stringify(body));
}

function setCorsHeaders(res) {
  res.setHeader('access-control-allow-origin', '*');
  res.setHeader('access-control-allow-methods', 'GET,POST,OPTIONS');
  res.setHeader('access-control-allow-headers', 'content-type,t,s');
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let body = '';

    req.setEncoding('utf8');
    req.on('data', (chunk) => {
      body += chunk;

      if (body.length > 1024 * 1024) {
        req.destroy();
        reject(httpError(413, 'Request body is too large'));
      }
    });
    req.on('end', () => resolve(body));
    req.on('error', reject);
  });
}

function httpError(statusCode, message) {
  const error = new Error(message);
  error.statusCode = statusCode;
  return error;
}

function trimTrailingSlash(value) {
  return value.replace(/\/+$/, '');
}

function loadDotEnv() {
  if (!fs.existsSync('.env')) return;

  const lines = fs.readFileSync('.env', 'utf8').split(/\r?\n/);

  for (const line of lines) {
    const trimmed = line.trim();

    if (!trimmed || trimmed.startsWith('#')) continue;

    const separatorIndex = trimmed.indexOf('=');
    if (separatorIndex === -1) continue;

    const key = trimmed.slice(0, separatorIndex).trim();
    const value = trimmed.slice(separatorIndex + 1).trim().replace(/^["']|["']$/g, '');

    if (key && process.env[key] === undefined) {
      process.env[key] = value;
    }
  }
}
