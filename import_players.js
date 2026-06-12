const mysql = require('mysql2/promise');

const DB_CONFIG = {
    host: '192.168.31.23',
    port: 13306,
    user: 'root',
    password: 'Sun123456',
    database: 'gva',
    charset: 'utf8mb4'
};

// 战队映射 (用户提供)
const teamMap = [
    { ID: 2, teamName: '狂飙' },
    { ID: 3, teamName: 'HS战队' },
    { ID: 4, teamName: 'ACY战队' },
    { ID: 5, teamName: 'LY战队' },
    { ID: 6, teamName: 'ZZL' },
    { ID: 7, teamName: '人坤' },
    { ID: 8, teamName: 'PP战队' },
    { ID: 9, teamName: 'T1' },
    { ID: 10, teamName: '中国风' },
    { ID: 11, teamName: '粉兔兔' },
    { ID: 12, teamName: 'ALL' },
    { ID: 13, teamName: '耙耳朵' },
    { ID: 14, teamName: 'TGH' },
    { ID: 15, teamName: 'Galaxy Fleet' },
    { ID: 16, teamName: '猛虎' },
    { ID: 17, teamName: '花露水' },
    { ID: 18, teamName: 'W战队' },
    { ID: 19, teamName: '夕阳红' },
    { ID: 20, teamName: '新势力' },
    { ID: 21, teamName: 'Myths' },
    { ID: 22, teamName: 'DFH' },
    { ID: 23, teamName: '彩虹岛战队' },
    { ID: 24, teamName: '权威' },
    { ID: 25, teamName: 'DH黑马战队' },
    { ID: 26, teamName: 'yy' },
    { ID: 27, teamName: '缔造' },
    { ID: 28, teamName: '极致' },
    { ID: 29, teamName: '哈哈' },
    { ID: 30, teamName: '逆境' },
    { ID: 31, teamName: 'SC' }
];

// 选手UID数据 (从Excel读取)
const playersData = [
    // 狂飙 (ID=2)
    { teamId: 2, uid: 'ziac2' },
    { teamId: 2, uid: '186aq9' },
    { teamId: 2, uid: '182ymp' },
    { teamId: 2, uid: 'zkkia' },
    { teamId: 2, uid: 'zizma' },
    // HS战队 (ID=3)
    { teamId: 3, uid: 'zchqq' },
    { teamId: 3, uid: '17ucwh' },
    { teamId: 3, uid: 'zgvrm' },
    { teamId: 3, uid: '17wgr7' },
    { teamId: 3, uid: 'zh236' },
    // ACY战队 (ID=4)
    { teamId: 4, uid: '17vxsh' },
    { teamId: 4, uid: 'zhxoy' },
    { teamId: 4, uid: '1834y9' },
    { teamId: 4, uid: '17x61f' },
    { teamId: 4, uid: '183hld' },
    // LY战队 (ID=5)
    { teamId: 5, uid: 'z6cia' },
    { teamId: 5, uid: 'zmuoi' },
    { teamId: 5, uid: '1858sx' },
    { teamId: 5, uid: 'zjiky' },
    { teamId: 5, uid: 'zkqtu' },
    // ZZL (ID=6)
    { teamId: 6, uid: '17rq37' },
    { teamId: 6, uid: 'zlg42' },
    { teamId: 6, uid: 'zn102' },
    { teamId: 6, uid: 'zddci' },
    { teamId: 6, uid: '184jip' },
    // 人坤 (ID=7)
    { teamId: 7, uid: '17xvbn' },
    { teamId: 7, uid: '17zsur' },
    { teamId: 7, uid: '184pu9' },
    { teamId: 7, uid: 'zk1jm' },
    { teamId: 7, uid: '17yx8z' },
    // PP战队 (ID=8)
    { teamId: 8, uid: '17yklv' },
    { teamId: 8, uid: 'zfavm' },
    { teamId: 8, uid: '17nopf' },
    { teamId: 8, uid: 'zbm4y' },
    { teamId: 8, uid: '183nwx' },
    // T1 (ID=9)
    { teamId: 9, uid: '17rdg1' },
    { teamId: 9, uid: '17t4nl' },
    { teamId: 9, uid: '17taz5' },
    { teamId: 9, uid: '17xccz' },
    { teamId: 9, uid: 'zimz6' },
    // 中国风 (ID=10)
    { teamId: 10, uid: '17thap' },
    { teamId: 10, uid: 'zawuq' },
    { teamId: 10, uid: '17zg7n' },
    { teamId: 10, uid: 'zl9si' },
    { teamId: 10, uid: 'zi40i' },
    // 粉兔兔 (ID=11)
    { teamId: 11, uid: 'z9iaa' },
    { teamId: 11, uid: 'zkx5e' },
    { teamId: 11, uid: 'zitaq' },
    { teamId: 11, uid: 'z95n6' },
    { teamId: 11, uid: '1846vl' },
    // ALL (ID=12)
    { teamId: 12, uid: '17wzpv' },
    { teamId: 12, uid: '17z9w3' },
    { teamId: 12, uid: 'zhrde' },
    { teamId: 12, uid: '17qo5v' },
    { teamId: 12, uid: 'znqaa' },
    // 耙耳朵 (ID=13)
    { teamId: 13, uid: '181k29' },
    { teamId: 13, uid: '181wpd' },
    { teamId: 13, uid: '17wteb' },
    { teamId: 13, uid: 'zheqa' },
    { teamId: 13, uid: '17wn2r' },
    // TGH (ID=14)
    { teamId: 14, uid: '17vxsj' },
    { teamId: 14, uid: 'zm5ea' },
    { teamId: 14, uid: '17xioj' },
    { teamId: 14, uid: '17vrgz' },
    { teamId: 14, uid: 'zhl1u' },
    // Galaxy Fleet (ID=15)
    { teamId: 15, uid: '17r74h' },
    { teamId: 15, uid: '17n5qr' },
    { teamId: 15, uid: '17w443' },
    { teamId: 15, uid: '18230x' },
    { teamId: 15, uid: '17xp03' },
    // 猛虎 (ID=16)
    { teamId: 16, uid: '17yeab' },
    { teamId: 16, uid: '17odzn' },
    { teamId: 16, uid: '17xp01' },
    { teamId: 16, uid: '17o7o3' },
    { teamId: 16, uid: '17w441' },
    // 花露水 (ID=17)
    { teamId: 17, uid: 'zc53m' },
    { teamId: 17, uid: 'z8t02' },
    { teamId: 17, uid: '17q577' },
    { teamId: 17, uid: 'zaqj6' },
    { teamId: 17, uid: '182sb5' },
    // W战队 (ID=18)
    { teamId: 18, uid: '1829ch' },
    { teamId: 18, uid: 'zj5xu' },
    { teamId: 18, uid: 'z71si' },
    { teamId: 18, uid: '17wafn' },
    { teamId: 18, uid: '180i4z' },
    // 夕阳红 (ID=19)
    { teamId: 19, uid: 'zignm' },
    { teamId: 19, uid: 'zgpg2' },
    { teamId: 19, uid: '17y1n7' },
    { teamId: 19, uid: '181qdt' },
    { teamId: 19, uid: 'zh8eq' },
    // 新势力 (ID=20)
    { teamId: 20, uid: 'zn7bm' },
    { teamId: 20, uid: '182lzl' },
    { teamId: 20, uid: '183b9t' },
    { teamId: 20, uid: '182fo1' },
    { teamId: 20, uid: '17rjrl' },
    // Myths (ID=21)
    { teamId: 21, uid: '17r74j' },
    { teamId: 21, uid: 'z6vgy' },
    { teamId: 21, uid: 'z7r2q' },
    { teamId: 21, uid: 'zf4k2' },
    { teamId: 21, uid: 'zd70y' },
    // DFH (ID=22)
    { teamId: 22, uid: '17v8i9' },
    { teamId: 22, uid: '185rrl' },
    { teamId: 22, uid: '17xioh' },
    { teamId: 22, uid: '180btf' },
    { teamId: 22, uid: '17u6kx' },
    // 彩虹岛战队 (ID=23)
    { teamId: 23, uid: 'zjc9e' },
    { teamId: 23, uid: 'zlsr6' },
    { teamId: 23, uid: '17wn2p' },
    { teamId: 23, uid: 'zjowi' },
    { teamId: 23, uid: 'z8moi' },
    // 权威 (ID=24)
    { teamId: 24, uid: 'za18y' },
    { teamId: 24, uid: 'zjv82' },
    { teamId: 24, uid: '17y7yr' },
    { teamId: 24, uid: '17owyb' },
    { teamId: 24, uid: '17wte9' },
    // DH黑马战队 (ID=25)
    { teamId: 25, uid: '17z3kj' },
    { teamId: 25, uid: '17uj81' },
    { teamId: 25, uid: '17psk3' },
    { teamId: 25, uid: 'zcudu' },
    { teamId: 25, uid: 'zco2a' },
    // yy (ID=26)
    { teamId: 26, uid: '17ttxt' },
    { teamId: 26, uid: 'zmbpu' },
    { teamId: 26, uid: '17zmj7' },
    { teamId: 26, uid: '1852hd' },
    { teamId: 26, uid: '184w5t' },
    // 缔造 (ID=27)
    { teamId: 27, uid: 'zk7v6' },
    { teamId: 27, uid: 'zke6q' },
    { teamId: 27, uid: '1840k1' },
    { teamId: 27, uid: '183u8h' },
    { teamId: 27, uid: '17qo5t' },
    // 极致 (ID=28)
    { teamId: 28, uid: 'zlmfm' },
    { teamId: 28, uid: 'zl3gy' },
    { teamId: 28, uid: '184d75' },
    { teamId: 28, uid: '180ogj' },
    { teamId: 28, uid: '17yqxf' },
    // 哈哈 (ID=29)
    { teamId: 29, uid: '17r0sz' },
    { teamId: 29, uid: 'zmi1e' },
    { teamId: 29, uid: '17zz6b' },
    { teamId: 29, uid: '17p39v' },
    { teamId: 29, uid: '17u09d' },
    // 逆境 (ID=30) - 只有4人
    { teamId: 30, uid: 'zlz2q' },
    { teamId: 30, uid: '185y35' },
    { teamId: 30, uid: '1864ep' },
    { teamId: 30, uid: '17z9w1' },
    // SC (ID=31)
    { teamId: 31, uid: '185f4h' },
    { teamId: 31, uid: 'zmocy' },
    { teamId: 31, uid: '185lg1' },
    { teamId: 31, uid: '17slox' },
    { teamId: 31, uid: 'z9uxe' }
];

// 创建战队ID到战队名的映射
const teamIdToName = new Map();
teamMap.forEach(t => teamIdToName.set(t.ID, t.teamName));

async function importPlayers() {
    let connection;
    const results = {
        success: true,
        imported: 0,
        skipped: 0,
        errors: []
    };

    try {
        console.log('========== 导入选手数据 ==========\n');

        connection = await mysql.createConnection(DB_CONFIG);
        console.log('数据库连接成功!\n');

        await connection.beginTransaction();

        // 按战队统计成员序号
        const teamMemberCount = new Map();
        for (const player of playersData) {
            const count = teamMemberCount.get(player.teamId) || 0;
            teamMemberCount.set(player.teamId, count + 1);
        }

        console.log('===== 开始导入选手 =====\n');

        for (const player of playersData) {
            const teamName = teamIdToName.get(player.teamId);
            if (!teamName) {
                results.errors.push({ uid: player.uid, error: `战队ID ${player.teamId} 不存在` });
                continue;
            }

            // 获取该战队当前成员数量，用于生成序号
            const currentCount = teamMemberCount.get(player.teamId);
            const playerName = `${teamName}_${currentCount}`;
            teamMemberCount.set(player.teamId, currentCount - 1);

            try {
                // 检查是否已存在
                const [existing] = await connection.execute(
                    'SELECT id FROM exa_players WHERE uid = ?',
                    [player.uid]
                );

                if (existing.length > 0) {
                    // 更新现有记录
                    await connection.execute(
                        'UPDATE exa_players SET player_name = ?, team_id = ? WHERE uid = ?',
                        [playerName, player.teamId, player.uid]
                    );
                    console.log(`  [更新] UID=${player.uid} -> "${playerName}" (战队: ${teamName})`);
                } else {
                    // 插入新记录
                    await connection.execute(
                        'INSERT INTO exa_players (player_name, uid, team_id, bounty) VALUES (?, ?, ?, ?)',
                        [playerName, player.uid, player.teamId, 0]
                    );
                    console.log(`  [新增] UID=${player.uid} -> "${playerName}" (战队: ${teamName})`);
                    results.imported++;
                }
            } catch (err) {
                console.log(`  [错误] UID=${player.uid}: ${err.message}`);
                results.errors.push({ uid: player.uid, error: err.message });
            }
        }

        await connection.commit();
        console.log('\n事务已提交!');

        // 验证数据
        console.log('\n===== 数据验证 =====');
        const [playerCount] = await connection.execute('SELECT COUNT(*) as count FROM exa_players');
        const [teamStats] = await connection.execute(`
            SELECT t.team_name, COUNT(p.id) as player_count
            FROM exa_teams t
            LEFT JOIN exa_players p ON t.id = p.team_id
            GROUP BY t.id, t.team_name
            ORDER BY t.id
        `);

        console.log(`选手总数: ${playerCount[0].count}`);
        console.log('\n各战队选手数量:');
        teamStats.forEach(t => {
            console.log(`  ${t.team_name}: ${t.player_count}人`);
        });

    } catch (err) {
        console.error('\n[错误]', err.message);
        results.success = false;
        if (connection) {
            await connection.rollback();
            console.log('事务已回滚');
        }
    } finally {
        if (connection) {
            await connection.end();
        }
    }

    console.log('\n========== 导入结果 ==========');
    console.log(`状态: ${results.success ? '成功' : '失败'}`);
    console.log(`新增: ${results.imported}`);
    console.log(`错误: ${results.errors.length}`);

    return results;
}

importPlayers();
