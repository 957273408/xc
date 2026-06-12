const mysql = require('mysql2/promise');
const XLSX = require('xlsx');
const path = require('path');

const DB_CONFIG = {
    host: '192.168.31.23',
    port: 13306,
    user: 'root',
    password: 'Sun123456',
    database: 'gva',
    charset: 'utf8mb4'
};

const filePath = path.join(__dirname, '战队信息表-2026-06-11 16_48_20(1).xlsx');

// 战队数据
const teams = [
    { id: 1, name: '狂飙', abbr: 'KB' },
    { id: 2, name: 'HS战队', abbr: 'HS' },
    { id: 3, name: 'ACY战队', abbr: 'ACY' },
    { id: 4, name: 'LY战队', abbr: 'LY' },
    { id: 5, name: 'ZZL', abbr: 'ZZL' },
    { id: 6, name: '人坤', abbr: 'RK' },
    { id: 7, name: 'PP战队', abbr: 'PP' },
    { id: 8, name: 'T1', abbr: 'T1' },
    { id: 9, name: '中国风', abbr: 'ZGF' },
    { id: 10, name: '粉兔兔', abbr: 'FTT' },
    { id: 11, name: 'ALL', abbr: 'ALL' },
    { id: 12, name: '耙耳朵', abbr: 'PED' },
    { id: 13, name: 'TGH', abbr: 'TGH' },
    { id: 14, name: 'Galaxy Fleet', abbr: 'GF' },
    { id: 15, name: '猛虎', abbr: 'MH' },
    { id: 16, name: '花露水', abbr: 'HLS' },
    { id: 17, name: 'W战队', abbr: 'W' },
    { id: 18, name: '夕阳红', abbr: 'XYH' },
    { id: 19, name: '新势力', abbr: 'NFS' },
    { id: 20, name: 'Myths', abbr: 'Myths' },
    { id: 21, name: 'DFH', abbr: 'DFH' },
    { id: 22, name: '彩虹岛战队', abbr: 'CHD' },
    { id: 23, name: '权威', abbr: 'qw' },
    { id: 24, name: 'DH黑马战队', abbr: 'DH' },
    { id: 25, name: 'yy', abbr: 'yy' },
    { id: 26, name: '缔造', abbr: 'Dz' },
    { id: 27, name: '极致', abbr: 'ACME' },
    { id: 28, name: '哈哈', abbr: 'HH' },
    { id: 29, name: '逆境', abbr: 'NJ' },
    { id: 30, name: 'SC', abbr: 'SC' }
];

// 选手数据 (从Excel读取)
const players = [
    { teamId: 1, teamName: '狂飙', uid: 'ziac2' },
    { teamId: 1, teamName: '狂飙', uid: '186aq9' },
    { teamId: 1, teamName: '狂飙', uid: '182ymp' },
    { teamId: 1, teamName: '狂飙', uid: 'zkkia' },
    { teamId: 1, teamName: '狂飙', uid: 'zizma' },
    { teamId: 2, teamName: 'HS战队', uid: 'zchqq' },
    { teamId: 2, teamName: 'HS战队', uid: '17ucwh' },
    { teamId: 2, teamName: 'HS战队', uid: 'zgvrm' },
    { teamId: 2, teamName: 'HS战队', uid: '17wgr7' },
    { teamId: 2, teamName: 'HS战队', uid: 'zh236' },
    { teamId: 3, teamName: 'ACY战队', uid: '17vxsh' },
    { teamId: 3, teamName: 'ACY战队', uid: 'zhxoy' },
    { teamId: 3, teamName: 'ACY战队', uid: '1834y9' },
    { teamId: 3, teamName: 'ACY战队', uid: '17x61f' },
    { teamId: 3, teamName: 'ACY战队', uid: '183hld' },
    { teamId: 4, teamName: 'LY战队', uid: 'z6cia' },
    { teamId: 4, teamName: 'LY战队', uid: 'zmuoi' },
    { teamId: 4, teamName: 'LY战队', uid: '1858sx' },
    { teamId: 4, teamName: 'LY战队', uid: 'zjiky' },
    { teamId: 4, teamName: 'LY战队', uid: 'zkqtu' },
    { teamId: 5, teamName: 'ZZL', uid: '17rq37' },
    { teamId: 5, teamName: 'ZZL', uid: 'zlg42' },
    { teamId: 5, teamName: 'ZZL', uid: 'zn102' },
    { teamId: 5, teamName: 'ZZL', uid: 'zddci' },
    { teamId: 5, teamName: 'ZZL', uid: '184jip' },
    { teamId: 6, teamName: '人坤', uid: '17xvbn' },
    { teamId: 6, teamName: '人坤', uid: '17zsur' },
    { teamId: 6, teamName: '人坤', uid: '184pu9' },
    { teamId: 6, teamName: '人坤', uid: 'zk1jm' },
    { teamId: 6, teamName: '人坤', uid: '17yx8z' },
    { teamId: 7, teamName: 'PP战队', uid: '17yklv' },
    { teamId: 7, teamName: 'PP战队', uid: 'zfavm' },
    { teamId: 7, teamName: 'PP战队', uid: '17nopf' },
    { teamId: 7, teamName: 'PP战队', uid: 'zbm4y' },
    { teamId: 7, teamName: 'PP战队', uid: '183nwx' },
    { teamId: 8, teamName: 'T1', uid: '17rdg1' },
    { teamId: 8, teamName: 'T1', uid: '17t4nl' },
    { teamId: 8, teamName: 'T1', uid: '17taz5' },
    { teamId: 8, teamName: 'T1', uid: '17xccz' },
    { teamId: 8, teamName: 'T1', uid: 'zimz6' },
    { teamId: 9, teamName: '中国风', uid: '17thap' },
    { teamId: 9, teamName: '中国风', uid: 'zawuq' },
    { teamId: 9, teamName: '中国风', uid: '17zg7n' },
    { teamId: 9, teamName: '中国风', uid: 'zl9si' },
    { teamId: 9, teamName: '中国风', uid: 'zi40i' },
    { teamId: 10, teamName: '粉兔兔', uid: 'z9iaa' },
    { teamId: 10, teamName: '粉兔兔', uid: 'zkx5e' },
    { teamId: 10, teamName: '粉兔兔', uid: 'zitaq' },
    { teamId: 10, teamName: '粉兔兔', uid: 'z95n6' },
    { teamId: 10, teamName: '粉兔兔', uid: '1846vl' },
    { teamId: 11, teamName: 'ALL', uid: '17wzpv' },
    { teamId: 11, teamName: 'ALL', uid: '17z9w3' },
    { teamId: 11, teamName: 'ALL', uid: 'zhrde' },
    { teamId: 11, teamName: 'ALL', uid: '17qo5v' },
    { teamId: 11, teamName: 'ALL', uid: 'znqaa' },
    { teamId: 12, teamName: '耙耳朵', uid: '181k29' },
    { teamId: 12, teamName: '耙耳朵', uid: '181wpd' },
    { teamId: 12, teamName: '耙耳朵', uid: '17wteb' },
    { teamId: 12, teamName: '耙耳朵', uid: 'zheqa' },
    { teamId: 12, teamName: '耙耳朵', uid: '17wn2r' },
    { teamId: 13, teamName: 'TGH', uid: '17vxsj' },
    { teamId: 13, teamName: 'TGH', uid: 'zm5ea' },
    { teamId: 13, teamName: 'TGH', uid: '17xioj' },
    { teamId: 13, teamName: 'TGH', uid: '17vrgz' },
    { teamId: 13, teamName: 'TGH', uid: 'zhl1u' },
    { teamId: 14, teamName: 'Galaxy Fleet', uid: '17r74h' },
    { teamId: 14, teamName: 'Galaxy Fleet', uid: '17n5qr' },
    { teamId: 14, teamName: 'Galaxy Fleet', uid: '17w443' },
    { teamId: 14, teamName: 'Galaxy Fleet', uid: '18230x' },
    { teamId: 14, teamName: 'Galaxy Fleet', uid: '17xp03' },
    { teamId: 15, teamName: '猛虎', uid: '17yeab' },
    { teamId: 15, teamName: '猛虎', uid: '17odzn' },
    { teamId: 15, teamName: '猛虎', uid: '17xp01' },
    { teamId: 15, teamName: '猛虎', uid: '17o7o3' },
    { teamId: 15, teamName: '猛虎', uid: '17w441' },
    { teamId: 16, teamName: '花露水', uid: 'zc53m' },
    { teamId: 16, teamName: '花露水', uid: 'z8t02' },
    { teamId: 16, teamName: '花露水', uid: '17q577' },
    { teamId: 16, teamName: '花露水', uid: 'zaqj6' },
    { teamId: 16, teamName: '花露水', uid: '182sb5' },
    { teamId: 17, teamName: 'W战队', uid: '1829ch' },
    { teamId: 17, teamName: 'W战队', uid: 'zj5xu' },
    { teamId: 17, teamName: 'W战队', uid: 'z71si' },
    { teamId: 17, teamName: 'W战队', uid: '17wafn' },
    { teamId: 17, teamName: 'W战队', uid: '180i4z' },
    { teamId: 18, teamName: '夕阳红', uid: 'zignm' },
    { teamId: 18, teamName: '夕阳红', uid: 'zgpg2' },
    { teamId: 18, teamName: '夕阳红', uid: '17y1n7' },
    { teamId: 18, teamName: '夕阳红', uid: '181qdt' },
    { teamId: 18, teamName: '夕阳红', uid: 'zh8eq' },
    { teamId: 19, teamName: '新势力', uid: 'zn7bm' },
    { teamId: 19, teamName: '新势力', uid: '182lzl' },
    { teamId: 19, teamName: '新势力', uid: '183b9t' },
    { teamId: 19, teamName: '新势力', uid: '182fo1' },
    { teamId: 19, teamName: '新势力', uid: '17rjrl' },
    { teamId: 20, teamName: 'Myths', uid: '17r74j' },
    { teamId: 20, teamName: 'Myths', uid: 'z6vgy' },
    { teamId: 20, teamName: 'Myths', uid: 'z7r2q' },
    { teamId: 20, teamName: 'Myths', uid: 'zf4k2' },
    { teamId: 20, teamName: 'Myths', uid: 'zd70y' },
    { teamId: 21, teamName: 'DFH', uid: '17v8i9' },
    { teamId: 21, teamName: 'DFH', uid: '185rrl' },
    { teamId: 21, teamName: 'DFH', uid: '17xioh' },
    { teamId: 21, teamName: 'DFH', uid: '180btf' },
    { teamId: 21, teamName: 'DFH', uid: '17u6kx' },
    { teamId: 22, teamName: '彩虹岛战队', uid: 'zjc9e' },
    { teamId: 22, teamName: '彩虹岛战队', uid: 'zlsr6' },
    { teamId: 22, teamName: '彩虹岛战队', uid: '17wn2p' },
    { teamId: 22, teamName: '彩虹岛战队', uid: 'zjowi' },
    { teamId: 22, teamName: '彩虹岛战队', uid: 'z8moi' },
    { teamId: 23, teamName: '权威', uid: 'za18y' },
    { teamId: 23, teamName: '权威', uid: 'zjv82' },
    { teamId: 23, teamName: '权威', uid: '17y7yr' },
    { teamId: 23, teamName: '权威', uid: '17owyb' },
    { teamId: 23, teamName: '权威', uid: '17wte9' },
    { teamId: 24, teamName: 'DH黑马战队', uid: '17z3kj' },
    { teamId: 24, teamName: 'DH黑马战队', uid: '17uj81' },
    { teamId: 24, teamName: 'DH黑马战队', uid: '17psk3' },
    { teamId: 24, teamName: 'DH黑马战队', uid: 'zcudu' },
    { teamId: 24, teamName: 'DH黑马战队', uid: 'zco2a' },
    { teamId: 25, teamName: 'yy', uid: '17ttxt' },
    { teamId: 25, teamName: 'yy', uid: 'zmbpu' },
    { teamId: 25, teamName: 'yy', uid: '17zmj7' },
    { teamId: 25, teamName: 'yy', uid: '1852hd' },
    { teamId: 25, teamName: 'yy', uid: '184w5t' },
    { teamId: 26, teamName: '缔造', uid: 'zk7v6' },
    { teamId: 26, teamName: '缔造', uid: 'zke6q' },
    { teamId: 26, teamName: '缔造', uid: '1840k1' },
    { teamId: 26, teamName: '缔造', uid: '183u8h' },
    { teamId: 26, teamName: '缔造', uid: '17qo5t' },
    { teamId: 27, teamName: '极致', uid: 'zlmfm' },
    { teamId: 27, teamName: '极致', uid: 'zl3gy' },
    { teamId: 27, teamName: '极致', uid: '184d75' },
    { teamId: 27, teamName: '极致', uid: '180ogj' },
    { teamId: 27, teamName: '极致', uid: '17yqxf' },
    { teamId: 28, teamName: '哈哈', uid: '17r0sz' },
    { teamId: 28, teamName: '哈哈', uid: 'zmi1e' },
    { teamId: 28, teamName: '哈哈', uid: '17zz6b' },
    { teamId: 28, teamName: '哈哈', uid: '17p39v' },
    { teamId: 28, teamName: '哈哈', uid: '17u09d' },
    { teamId: 29, teamName: '逆境', uid: 'zlz2q' },
    { teamId: 29, teamName: '逆境', uid: '185y35' },
    { teamId: 29, teamName: '逆境', uid: '1864ep' },
    { teamId: 29, teamName: '逆境', uid: '17z9w1' },
    // 逆境缺少第5名成员
    { teamId: 30, teamName: 'SC', uid: '185f4h' },
    { teamId: 30, teamName: 'SC', uid: 'zmocy' },
    { teamId: 30, teamName: 'SC', uid: '185lg1' },
    { teamId: 30, teamName: 'SC', uid: '17slox' },
    { teamId: 30, teamName: 'SC', uid: 'z9uxe' }
];

async function importData() {
    let connection;
    const results = {
        success: true,
        teamsImported: 0,
        teamsSkipped: 0,
        playersImported: 0,
        playersSkipped: 0,
        errors: []
    };

    try {
        console.log('========== 开始数据导入 ==========\n');
        console.log('正在连接数据库...');

        connection = await mysql.createConnection(DB_CONFIG);
        console.log('数据库连接成功!\n');

        // 开始事务
        await connection.beginTransaction();
        console.log('事务已启动\n');

        // ===== 1. 导入战队数据 =====
        console.log('===== 导入战队数据 =====');
        for (const team of teams) {
            try {
                // 检查战队是否已存在
                const [existing] = await connection.execute(
                    'SELECT id FROM exa_teams WHERE team_name = ?',
                    [team.name]
                );

                if (existing.length > 0) {
                    console.log(`  [跳过] 战队 "${team.name}" 已存在 (ID: ${existing[0].id})`);
                    results.teamsSkipped++;
                } else {
                    await connection.execute(
                        'INSERT INTO exa_teams (team_name, total_bounty) VALUES (?, ?)',
                        [team.name, 0]
                    );
                    console.log(`  [成功] 战队 "${team.name}" 已导入`);
                    results.teamsImported++;
                }
            } catch (err) {
                console.log(`  [错误] 战队 "${team.name}": ${err.message}`);
                results.errors.push({ type: 'team', name: team.name, error: err.message });
            }
        }
        console.log('');

        // ===== 2. 获取战队ID映射 =====
        console.log('===== 获取战队ID映射 =====');
        const [allTeams] = await connection.execute('SELECT id, team_name FROM exa_teams');
        const teamIdMap = new Map();
        for (const t of allTeams) {
            teamIdMap.set(t.team_name, t.id);
        }
        console.log(`已获取 ${teamIdMap.size} 个战队的ID映射\n`);

        // ===== 3. 导入选手数据 =====
        console.log('===== 导入选手数据 =====');
        for (const player of players) {
            try {
                const dbTeamId = teamIdMap.get(player.teamName);
                if (!dbTeamId) {
                    console.log(`  [错误] 无法找到战队 "${player.teamName}" 的ID`);
                    results.errors.push({ type: 'player', uid: player.uid, error: `Team not found: ${player.teamName}` });
                    continue;
                }

                // 检查选手是否已存在 (通过UID)
                const [existing] = await connection.execute(
                    'SELECT id FROM exa_players WHERE uid = ?',
                    [player.uid]
                );

                if (existing.length > 0) {
                    console.log(`  [跳过] 选手 UID="${player.uid}" (战队: ${player.teamName}) 已存在`);
                    results.playersSkipped++;
                } else {
                    // 选手姓名使用UID代替 (因为Excel中没有提供选手姓名)
                    await connection.execute(
                        'INSERT INTO exa_players (player_name, uid, team_id, bounty) VALUES (?, ?, ?, ?)',
                        [player.uid, player.uid, dbTeamId, 0]
                    );
                    console.log(`  [成功] 选手 UID="${player.uid}" (战队: ${player.teamName}) 已导入`);
                    results.playersImported++;
                }
            } catch (err) {
                console.log(`  [错误] 选手 UID="${player.uid}": ${err.message}`);
                results.errors.push({ type: 'player', uid: player.uid, error: err.message });
            }
        }

        // 提交事务
        await connection.commit();
        console.log('\n事务已提交!');

        // ===== 4. 验证数据 =====
        console.log('\n===== 数据验证 =====');
        const [teamCount] = await connection.execute('SELECT COUNT(*) as count FROM exa_teams');
        const [playerCount] = await connection.execute('SELECT COUNT(*) as count FROM exa_players');
        console.log(`  战队总数: ${teamCount[0].count}`);
        console.log(`  选手总数: ${playerCount[0].count}`);

    } catch (err) {
        console.error('\n[严重错误] 数据导入失败:', err.message);
        results.success = false;
        results.errors.push({ type: 'transaction', error: err.message });

        // 回滚事务
        if (connection) {
            await connection.rollback();
            console.log('事务已回滚');
        }
    } finally {
        if (connection) {
            await connection.end();
            console.log('数据库连接已关闭');
        }
    }

    // ===== 5. 输出结果汇总 =====
    console.log('\n========== 导入结果汇总 ==========');
    console.log(`操作状态: ${results.success ? '成功' : '失败'}`);
    console.log(`战队导入: ${results.teamsImported} 成功, ${results.teamsSkipped} 跳过`);
    console.log(`选手导入: ${results.playersImported} 成功, ${results.playersSkipped} 跳过`);
    console.log(`错误数量: ${results.errors.length}`);

    if (results.errors.length > 0) {
        console.log('\n错误详情:');
        results.errors.forEach((e, i) => {
            console.log(`  ${i + 1}. [${e.type}] ${e.name || e.uid || 'N/A'}: ${e.error}`);
        });
    }

    console.log('\n========== 数据导入完成 ==========');
    return results;
}

importData().then(results => {
    process.exit(results.success && results.errors.length === 0 ? 0 : 1);
}).catch(err => {
    console.error('未处理的错误:', err);
    process.exit(1);
});
