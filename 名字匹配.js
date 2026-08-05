/* ============================================================
 * 香肠派对 赛事观察台 - 接口联调模块
 * 仅 /xdc/get_info 和 /xdc/get_kill_info 可用（/xd/* 未部署，返回 404）
 * 积分需现场计算：score = rankScore(rank) + killScore * kills
 * ============================================================ */

(function (global) {
  'use strict';

  // ---------- 默认配置 ----------
  const DEFAULTS = {
    apiBase: 'http://sudcore.asia:18080/api/proxy',
    warId: '0fe848a9b617b6e3f0ccd5fe73a5594d_1',
    interval: 1000,
    killScore: 1,            // 每淘汰 1 分
    rankScoreTable: {        // PUBG Mobile 风格排名分
      1: 12, 2: 9, 3: 8, 4: 6, 5: 5, 6: 4, 7: 3, 8: 2,
      9: 1, 10: 1, 11: 1, 12: 1, 13: 0, 14: 0, 15: 0, 16: 0
    }
  };

  // ---------- 状态 ----------
  const state = {
    config: { ...DEFAULTS },
    timer: null,
    lastInfo: null,
    lastKills: null,
    lastError: null,
    lastSyncAt: 0,
    pollCount: 0,
    listeners: { data: [], error: [], status: [] }
  };

  // ---------- 工具函数 ----------
  const toArr = (v) => Array.isArray(v) ? v : (v ? [v] : []);

  const fmtNum = (v) => {
    if (v === null || v === undefined || v === '') return '-';
    const n = Number(v);
    if (isNaN(n)) return String(v);
    if (Number.isInteger(n)) return String(n);
    return n.toFixed(1);
  };

  // 从昵称中提取队伍前缀: "XHZ.辣子鸡" → "XHZ"; "XYT·繁仙" → "XYT"; "MS•兮辞" → "MS"; "STY．南茶" → "STY"
  // 支持的分隔符: 半角. 全角． 中点· 项目符号• 日文中点・ _ - # 空白
  const SEP = String.raw`[.\uFF0E\u00B7\u2022\u30FB_\-#\s]`;
  const getTeamFromNick = (nick) => {
    if (!nick) return '';
    const re = new RegExp('^([^' + SEP.slice(1, -1) + ']+)' + SEP);
    const m = String(nick).match(re);
    return m ? m[1] : '';
  };

  // 玩家状态: liveState 0存活/1倒地/2死亡; state 0正常/1倒地/2不可复活/3可复活; outZone 安全区外
  const getPlayerState = (p) => {
    if (!p) return 'alive';
    if (p.liveState === 2 || p.state === 2) return 'dead';      // 淘汰
    if (p.liveState === 1 || p.state === 1) return 'down';      // 倒地
    if (p.state === 3) return 'revivable';                       // 可复活
    if (p.outZone === true || p.inSafeZone === false) return 'out'; // 安全区外
    return 'alive';                                               // 存活
  };

  const rankScore = (rank) => {
    const r = Number(rank);
    if (!r) return 0;
    return state.config.rankScoreTable[r] ?? 0;
  };

  const calcScore = (rank, kills) => {
    return rankScore(rank) + (Number(kills) || 0) * state.config.killScore;
  };

  // ---------- 接口调用 ----------
  const buildUrl = (path, params) => {
    const base = (state.config.apiBase || '').replace(/\/+$/, '');
    const qs = params
      ? Object.entries(params)
          .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
          .join('&')
      : '';
    return `${base}${path}${qs ? '?' + qs : ''}`;
  };

  const request = async (path, params) => {
    const url = buildUrl(path, params);
    const resp = await fetch(url);
    if (!resp.ok) throw new Error(`HTTP ${resp.status} ${url}`);
    return resp.json();
  };

  const getInfo = (warId) => request('/xdc/get_info', { warId: warId || state.config.warId });
  const getKillInfo = (warId) => request('/xdc/get_kill_info', { warId: warId || state.config.warId });

  // ---------- 数据聚合：按队伍 ----------
  // 返回 [{ rank, teamTag, teamSize, aliveCount, downCount, deadCount, kills, score, members:[] }] 按 rank 升序
  const aggregateByTeam = (info) => {
    const players = toArr(info && info.playerInfoList);
    const teams = new Map();
    for (const p of players) {
      const tag = p.teamID || getTeamFromNick(p.nickName) || ('UID_' + p.uID);
      if (!teams.has(tag)) teams.set(tag, []);
      teams.get(tag).push(p);
    }
    const rows = [];
    for (const [tag, members] of teams) {
      // 队伍排名：取任一成员的 lostRoleRank（同队应一致）
      let rank = 0;
      for (const m of members) {
        if (m.lostRoleRank !== undefined && m.lostRoleRank !== null) {
          rank = Number(m.lostRoleRank) || 0;
          break;
        }
      }
      const aliveCount = members.filter(m => getPlayerState(m) === 'alive').length;
      const downCount = members.filter(m => getPlayerState(m) === 'down').length;
      const deadCount = members.filter(m => getPlayerState(m) === 'dead').length;
      const kills = members.reduce((s, m) => s + (Number(m.totalKill) || 0), 0);
      const score = calcScore(rank, kills);
      rows.push({
        rank, teamTag: tag, teamSize: members.length,
        aliveCount, downCount, deadCount, kills, score,
        members
      });
    }
    // 排序：rank 升序（rank=0 视为大数，置于末尾）；同 rank 时 kills 降序
    rows.sort((a, b) => {
      const ra = a.rank || 9999, rb = b.rank || 9999;
      if (ra !== rb) return ra - rb;
      return b.kills - a.kills;
    });
    // 若 rank 重复（同队不同成员可能值不一），重新分配序号 = 排名索引+1
    rows.forEach((r, i) => { r.displayRank = i + 1; });
    return rows;
  };

  // ---------- 状态分发 ----------
  function emit(event, payload) {
    const arr = state.listeners[event] || [];
    for (const cb of arr) {
      try { cb(payload); } catch (e) { console.error('[API] listener error', e); }
    }
  }

  function setStatus(status, text) {
    emit('status', { status, text });
  }

  // ---------- 单次拉取 ----------
  async function tick() {
    setStatus('loading', '查询中...');
    try {
      const [info, kills] = await Promise.all([
        getInfo(),
        getKillInfo().catch(e => null) // 击杀信息失败不阻塞主流程
      ]);
      state.lastInfo = info;
      state.lastKills = kills;
      state.lastSyncAt = Date.now();
      state.pollCount++;
      state.lastError = null;

      const players = toArr(info && info.playerInfoList);
      const teams = aggregateByTeam(info);
      emit('data', { info, kills, teams, players });
      setStatus('online', `已更新 ${new Date().toLocaleTimeString()}`);
    } catch (e) {
      state.lastError = e;
      setStatus('offline', '查询失败');
      emit('error', { message: e.message, error: e });
      console.error('[API] tick failed', e);
    }
  }

  // ---------- 轮询 ----------
  function start() {
    stop();
    tick(); // 立即触发一次
    state.timer = setInterval(tick, state.config.interval);
  }

  function stop() {
    if (state.timer) {
      clearInterval(state.timer);
      state.timer = null;
    }
  }

  // ---------- 配置 ----------
  function configure(cfg) {
    if (!cfg) return;
    Object.assign(state.config, cfg);
  }

  // ---------- 事件订阅 ----------
  function on(event, cb) {
    if (!state.listeners[event]) state.listeners[event] = [];
    state.listeners[event].push(cb);
    return () => {
      const arr = state.listeners[event];
      const i = arr.indexOf(cb);
      if (i >= 0) arr.splice(i, 1);
    };
  }

  // ---------- 导出 ----------
  const API = {
    state,
    configure,
    start,
    stop,
    tick,
    on,
    aggregateByTeam,
    toArr, fmtNum, getTeamFromNick, getPlayerState, rankScore, calcScore
  };

  global.SausageAPI = API;
})(window);
