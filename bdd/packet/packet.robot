*** Settings ***
Library           RequestsLibrary
Library           Collections
Library           String

*** Variables ***
${tokenHolder}    ${EMPTY}
${tokenHolderPubKey}    ${EMPTY}
${signature}      ${EMPTY}

*** Test Cases ***
packet
    [Documentation]    amount = 90
    ...    count = 10
    ...    min = 1
    ...    max = 10
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    sleep    3
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    90    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    4
    sleep    3
    pullPacket    ${tokenHolder}    4    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    5
    sleep    3
    pullPacket    ${tokenHolder}    5    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    6
    sleep    3
    pullPacket    ${tokenHolder}    6    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    7
    sleep    3
    pullPacket    ${tokenHolder}    7    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    8
    sleep    3
    pullPacket    ${tokenHolder}    8    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    9
    sleep    3
    pullPacket    ${tokenHolder}    9    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    10
    sleep    3
    pullPacket    ${tokenHolder}    10    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    90
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    0
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    10    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    90
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    10
    Should Be Equal As Strings    ${pulled}    true
    sleep    3
    sign    ${twoAddr}    11
    sleep    3
    pullPacket    ${tokenHolder}    11    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    90
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}

packet1
    [Documentation]    amount = 90
    ...    count = 10
    ...    min = 1
    ...    max = 10
    ...
    ...    调整为
    ...
    ...    amount = 90
    ...    count = 11
    ...    min = 2
    ...    max = 11
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    90    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    updatePacket    ${twoAddr}    ${twoAddr}    0    ${tokenHolderPubKey}    11    2
    ...    11    ${EMPTY}    false
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceCount"]}    11
    getBalance    ${twoAddr}    PTN
    getAllPacketInfo

packet2
    [Documentation]    amount = 90
    ...    count = 10
    ...    min = 1
    ...    max = 10
    ...
    ...    调整为
    ...
    ...    amount = 100
    ...    count = 11
    ...    min = 2
    ...    max = 11
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    90    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    updatePacket    ${twoAddr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    10    ${tokenHolderPubKey}    11    2
    ...    11    ${EMPTY}    TRUE
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    100

packet3
    [Documentation]    amount = 9
    ...    count = 10
    ...    min = 1
    ...    max = 10
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    9    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    4
    sleep    3
    pullPacket    ${tokenHolder}    4    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    5
    sleep    3
    pullPacket    ${tokenHolder}    5    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    6
    sleep    3
    pullPacket    ${tokenHolder}    6    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    7
    sleep    3
    pullPacket    ${tokenHolder}    7    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    8
    sleep    3
    pullPacket    ${tokenHolder}    8    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    9
    sleep    3
    pullPacket    ${tokenHolder}    9    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    0
    Should Be Equal As Strings    ${result["BalanceCount"]}    1
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    9    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    9
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    9
    Should Be Equal As Strings    ${pulled}    true
    sleep    3
    sign    ${twoAddr}    10
    sleep    3
    pullPacket    ${tokenHolder}    10    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}

packet4
    [Documentation]    amount = 900
    ...    count = 10
    ...    min = 1
    ...    max = 10
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    900    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    true
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    11
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    1
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    22
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    2
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    33
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    6
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    894
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    3
    Should Be Equal As Strings    ${pulled}    true

packet5
    [Documentation]    红包过期退回
    ${time}    Get Time
    log    ${time}
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    sleep    3
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    900    ${tokenHolderPubKey}    10    1    10
    ...    2020-05-06 14:51:18    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    getBalance    ${twoAddr}    PTN
    recyclePacket    ${twoAddr}    ${tokenHolderPubKey}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    0
    Should Be Equal As Strings    ${result["BalanceCount"]}    0
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${amount}    getBalance    ${twoAddr}    PTN
    Should Be Equal As Numbers    ${amount}    9998
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${tokenHolder}    0

packet6
    [Documentation]    amount = 900
    ...    count = 10
    ...    min = 1
    ...    max = 10
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    900    ${tokenHolderPubKey}    10    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    4
    sleep    3
    pullPacket    ${tokenHolder}    4    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    5
    sleep    3
    pullPacket    ${tokenHolder}    5    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    6
    sleep    3
    pullPacket    ${tokenHolder}    6    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    7
    sleep    3
    pullPacket    ${tokenHolder}    7    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    8
    sleep    3
    pullPacket    ${tokenHolder}    8    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    9
    sleep    3
    pullPacket    ${tokenHolder}    9    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    10
    sleep    3
    pullPacket    ${tokenHolder}    10    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    100
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    800
    Should Be Equal As Strings    ${result["BalanceCount"]}    0
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sleep    3
    pullPacket    ${tokenHolder}    10    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    10
    Should Be Equal As Strings    ${pulled}    true
    sleep    3
    sign    ${twoAddr}    11
    sleep    3
    pullPacket    ${tokenHolder}    11    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}

packet7
    [Documentation]    amount = 30
    ...    count = 0
    ...    min = 1
    ...    max = 10
    ...
    ...    count = 0====》无限领取
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    sleep    3
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    30    ${tokenHolderPubKey}    0    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    0
    Should Be Equal As Strings    ${result["BalanceCount"]}    0
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    4
    sleep    3
    pullPacket    ${tokenHolder}    4    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    30
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    3
    Should Be Equal As Strings    ${pulled}    true
    getAllPacketInfo

packet8
    [Documentation]    amount = 30
    ...    count = 3
    ...    min = 1
    ...    max = 10
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    sleep    3
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    sleep    3
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    createPacket    ${twoAddr}    30    ${tokenHolderPubKey}    3    1    10
    ...    ${EMPTY}    false
    sleep    3
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["BalanceAmount"]}    0
    Should Be Equal As Strings    ${result["BalanceCount"]}    0
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    4
    sleep    3
    pullPacket    ${tokenHolder}    4    ${signature}    ${oneAddr}    0
    sleep    3
    getBalance    ${oneAddr}    PTN
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    0
    sleep    3
    ${amount}    getBalance    ${oneAddr}    PTN
    Should Be Equal As Numbers    ${amount}    30
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    isPulledPacket    ${tokenHolderPubKey}    3
    getAllPacketInfo

multiToken
    [Documentation]    创建包含 t1 和 t2 多 token 的红包，并且领取 3 次
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount    #获取红包测试地址
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}    #获取红包测试地址的公钥
    ${result}    createToken    ${twoAddr}    t1
    log    ${result}
    sleep    5
    ${assetId1}    ccquery    t1
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    log    ${t}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    PTN
    log    ${t}
    ${result}    createToken    ${twoAddr}    t2
    log    ${result}
    sleep    5
    ${assetId2}    ccquery    t2
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000
    createMultiTokenPacket    ${twoAddr}    90    ${tokenHolderPubKey}    0    0    0
    ...    ${EMPTY}    true    ${assetId1}    ${assetId2}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    90
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    90
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    11,1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    1,1
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    22,2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    2,2
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    33,3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3,3
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3,3
    sleep    3
    ${amount}    getBalance    ${oneAddr}    ${assetId1}
    Should Be Equal As Numbers    ${amount}    6
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    84
    log    ${result}
    ${result}    getPacketAllocationHistory    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result[0]["PubKey"]}    ${tokenHolderPubKey}
    ${pulled}    isPulledPacket    ${tokenHolderPubKey}    3
    Should Be Equal As Strings    ${pulled}    true

multiTokenUpdate
    [Documentation]    t3 token 从
    ...
    ...
    ...    amount = 90
    ...    count = 0
    ...    min = 0
    ...    max = 0
    ...
    ...    调整为
    ...
    ...    amount = 100
    ...    count = 0
    ...    min = 0
    ...    max = 0
    ...
    ...
    ...    t4 token 从
    ...
    ...
    ...    amount = 90
    ...    count = 0
    ...    min = 0
    ...    max = 0
    ...
    ...    调整为
    ...
    ...    amount = 100
    ...    count = 0
    ...    min = 0
    ...    max = 0
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    ${result}    createToken    ${twoAddr}    t3
    log    ${result}
    sleep    5
    ${assetId1}    ccquery    t3
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    log    ${t}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    PTN
    log    ${t}
    ${result}    createToken    ${twoAddr}    t4
    log    ${result}
    sleep    5
    ${assetId2}    ccquery    t4
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000
    createMultiTokenPacket    ${twoAddr}    90    ${tokenHolderPubKey}    0    0    0
    ...    ${EMPTY}    true    ${assetId1}    ${assetId2}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    90
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    90
    multiTokenUpdated    ${twoAddr}    ${assetId1}    10    ${assetId2}    10
    sleep    4
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    100
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    100

multiAppend
    [Documentation]    创建包含 t5 和 t6 多 token 的红包，并且最后追加 PTN token 进入
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    ${result}    createToken    ${twoAddr}    t5
    log    ${result}
    sleep    5
    ${assetId1}    ccquery    t5
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    log    ${t}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    PTN
    log    ${t}
    ${result}    createToken    ${twoAddr}    t6
    log    ${result}
    sleep    5
    ${assetId2}    ccquery    t6
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000
    createMultiTokenPacket    ${twoAddr}    90    ${tokenHolderPubKey}    0    0    0
    ...    ${EMPTY}    true    ${assetId1}    ${assetId2}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    90
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    90
    updatePacket    ${twoAddr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    100    ${tokenHolderPubKey}    0    0
    ...    0    ${EMPTY}    true
    sleep    4
    multiTokenUpdated    ${twoAddr}    ${assetId1}    10    ${assetId2}    10
    sleep    4
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    100
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    100
    Should Be Equal As Strings    ${result["Token"][2]["amount"]}    100

multiTokenRecycle
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}
    ${result}    createToken    ${twoAddr}    t7
    log    ${result}
    sleep    3
    ${assetId1}    ccquery    t7
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    log    ${t}
    Should Be Equal As Numbers    ${t}    1000
    ${result}    createToken    ${twoAddr}    t8
    log    ${result}
    sleep    3
    ${assetId2}    ccquery    t8
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    PTN
    log    ${t}
    Should Be Equal As Numbers    ${t}    9998
    createMultiTokenPacket    ${twoAddr}    90    ${tokenHolderPubKey}    0    0    0
    ...    2020-05-06 14:51:18    true    ${assetId1}    ${assetId2}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    90
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    90
    recyclePacket    ${twoAddr}    ${tokenHolderPubKey}
    sleep    5
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    0
    Should Be Equal As Strings    ${result["Token"][1]["BalanceAmount"]}    0
    ${t}    getBalance    ${twoAddr}    PTN
    Should Be Equal As Numbers    ${t}    9996
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000

multiTokenPull
    [Documentation]    领完不够，更新
    listAccounts    #    主要获取 tokenHolder
    unlockAccount    ${tokenHolder}    1    #    解锁 tokenHolder
    ${twoAddr}    newAccount    #获取红包测试地址
    sleep    3
    transferPtn    ${tokenHolder}    ${twoAddr}    10000    1    1
    sleep    3
    unlockAccount    ${twoAddr}    1
    getBalance    ${twoAddr}    PTN
    getPublicKey    ${twoAddr}    #获取红包测试地址的公钥
    ${result}    createToken    ${twoAddr}    t9
    log    ${result}
    sleep    5
    ${assetId1}    ccquery    t9
    ${t}    getBalance    ${twoAddr}    ${assetId1}
    log    ${t}
    Should Be Equal As Numbers    ${t}    1000
    ${t}    getBalance    ${twoAddr}    PTN
    log    ${t}
    ${result}    createToken    ${twoAddr}    t10
    log    ${result}
    sleep    5
    ${assetId2}    ccquery    t10
    ${t}    getBalance    ${twoAddr}    ${assetId2}
    Should Be Equal As Numbers    ${t}    1000
    createMultiTokenPacket    ${twoAddr}    5    ${tokenHolderPubKey}    0    0    0
    ...    ${EMPTY}    true    ${assetId1}    ${assetId2}
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    5
    Should Be Equal As Strings    ${result["Token"][1]["amount"]}    5
    getPacketAllocationHistory    ${tokenHolderPubKey}
    ${oneAddr}    newAccount
    sleep    3
    sign    ${twoAddr}    11,1
    sleep    3
    pullPacket    ${tokenHolder}    1    ${signature}    ${oneAddr}    1,1
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    22,2
    sleep    3
    pullPacket    ${tokenHolder}    2    ${signature}    ${oneAddr}    2,2
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    sign    ${twoAddr}    33,3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3,3    # 不成功
    sleep    3
    getBalance    ${oneAddr}    ${assetId1}
    getPacketInfo    ${tokenHolderPubKey}
    getPacketAllocationHistory    ${tokenHolderPubKey}
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3,3    # 不成功
    sleep    3
    ${amount}    getBalance    ${oneAddr}    ${assetId1}
    Should Be Equal As Numbers    ${amount}    3
    ${amount}    getBalance    ${oneAddr}    ${assetId2}
    Should Be Equal As Numbers    ${amount}    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    2
    Should Be Equal As Strings    ${result["Token"][1]["BalanceAmount"]}    2
    ${result}    getPacketAllocationHistory    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result[0]["PubKey"]}    ${tokenHolderPubKey}
    multiTokenUpdated    ${twoAddr}    ${assetId1}    10    ${assetId2}    10
    sleep    3
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    12
    Should Be Equal As Strings    ${result["Token"][1]["BalanceAmount"]}    12
    sign    ${twoAddr}    33,3
    sleep    3
    pullPacket    ${tokenHolder}    3    ${signature}    ${oneAddr}    3,3    # 成功
    sleep    3
    ${amount}    getBalance    ${oneAddr}    ${assetId1}
    Should Be Equal As Numbers    ${amount}    6
    ${amount}    getBalance    ${oneAddr}    ${assetId2}
    Should Be Equal As Numbers    ${amount}    6
    ${result}    getPacketInfo    ${tokenHolderPubKey}
    Should Be Equal As Strings    ${result["Token"][0]["BalanceAmount"]}    9
    Should Be Equal As Strings    ${result["Token"][1]["BalanceAmount"]}    9

*** Keywords ***
createPacket
    [Arguments]    ${addr}    ${amount}    ${pubkey}    ${count}    ${min}    ${max}
    ...    ${expiredTime}    ${isConstant}
    ${param}    Create List    createPacket    ${pubkey}    ${count}    ${min}    ${max}
    ...    ${expiredTime}    remark    ${isConstant}
    ${two}    Create List    ${addr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${amount}    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99
    ...    ${param}
    ${res}    post    contract_ccinvoketx    createPacket    ${two}
    log    ${res}    #    #    Create List    createPacket    ${pubkey}
    ...    # ${count}    ${min}    ${max}    # ${expiredTime}    remark    #
    ...    # Create List    ${addr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    PTN    ${amount}    1
    ...    # PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}
    #    post    contract_ccinvokeToken    createPacket    ${two}
    #    ${res}
    [Return]    ${res}

post
    [Arguments]    ${method}    ${alias}    ${params}
    ${header}    Create Dictionary    Content-Type=application/json
    ${data}    Create Dictionary    jsonrpc=2.0    method=${method}    params=${params}    id=1
    Create Session    ${alias}    http://127.0.0.1:8545    #    http://127.0.0.1:8545    http://192.168.44.128:8545
    ${resp}    Post Request    ${alias}    http://127.0.0.1:8545    data=${data}    headers=${header}
    ${respJson}    To Json    ${resp.content}
    Dictionary Should Contain Key    ${respJson}    result
    ${res}    Get From Dictionary    ${respJson}    result
    [Return]    ${res}

listAccounts
    ${param}    Create List
    ${result}    post    personal_listAccounts    personal_listAccounts    ${param}
    log    ${result}
    Set Global Variable    ${tokenHolder}    ${result[0]}
    log    ${tokenHolder}

unlockAccount
    [Arguments]    ${addr}    ${pwd}
    ${param}    Create List    ${addr}    ${pwd}
    ${result}    post    personal_unlockAccount    personal_unlockAccount    ${param}
    log    ${result}
    Should Be True    ${result}

getPublicKey
    [Arguments]    ${addr}
    ${param}    Create List    ${addr}    1
    ${result}    post    personal_getPublicKey    personal_getPublicKey    ${param}
    log    ${result}
    Set Global Variable    ${tokenHolderPubKey}    ${result}
    log    ${tokenHolderPubKey}

getPacketInfo
    [Arguments]    ${tokenHolderPubKey}
    ${param}    Create List    getPacketInfo    ${tokenHolderPubKey}
    ${two}    Create List    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}    ${10}
    ${res}    post    contract_ccquery    getPacketInfo    ${two}
    log    ${res}
    ${addressMap}    To Json    ${res}
    [Return]    ${addressMap}

pullPacket
    [Arguments]    ${addr}    ${message}    ${signature}    ${pullAddr}    ${amount}
    ${param}    Create List    pullPacket    ${tokenHolderPubKey}    ${message}    ${signature}    ${pullAddr}
    ...    ${amount}
    ${two}    Create List    ${addr}    ${addr}    1    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99
    ...    ${param}
    ${res}    post    contract_ccinvoketx    pullPacket    ${two}
    log    ${res}

getBalance
    [Arguments]    ${address}    ${assetId}
    ${two}    Create List    ${address}
    ${result}    post    wallet_getBalance    wallet_getBalance    ${two}
    log    ${result}
    ${len}    Get Length    ${result}
    ${amount}    Set Variable If    ${len} == 0    0    ${result["${assetId}"]}
    [Return]    ${amount}

getPacketAllocationHistory
    [Arguments]    ${pubkey}
    ${param}    Create List    getPacketAllocationHistory    ${pubkey}
    ${two}    Create List    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}    ${10}
    ${res}    post    contract_ccquery    getPacketAllocationHistory    ${two}
    log    ${res}
    ${addressMap}    To Json    ${res}
    [Return]    ${addressMap}

updatePacket
    [Arguments]    ${addr}    ${toaddr}    ${amount}    ${pubkey}    ${count}    ${min}
    ...    ${max}    ${expiredTime}    ${isConstant}
    ${param}    Create List    updatePacket    ${pubkey}    ${count}    ${min}    ${max}
    ...    ${expiredTime}    remark    ${isConstant}
    ${two}    Create List    ${addr}    ${toaddr}    ${amount}    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99
    ...    ${param}
    ${res}    post    contract_ccinvoketx    updatePacket    ${two}
    log    ${res}
    [Return]    ${res}

recyclePacket
    [Arguments]    ${addr}    ${pubkey}
    ${param}    Create List    recyclePacket    ${pubkey}
    ${two}    Create List    ${addr}    ${addr}    0    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99
    ...    ${param}
    ${res}    post    contract_ccinvoketx    recyclePacket    ${two}
    log    ${res}

newAccount
    ${param}    Create List    1
    ${result}    post    personal_newAccount    personal_newAccount    ${param}
    log    ${result}
    #    ${oneAddr}    ${result}
    [Return]    ${result}

transferPtn
    [Arguments]    ${fromAddr}    ${toAddr}    ${amount}    ${fee}    ${pwd}
    ${param}    Create List    ${fromAddr}    ${toAddr}    ${amount}    ${fee}    ${null}
    ...    ${pwd}
    ${result}    post    wallet_transferPtn    wallet_transferPtn    ${param}
    log    ${result}

sign
    [Arguments]    ${addr}    ${message}
    ${param}    Create List    ${message}    ${addr}    1
    ${result}    post    personal_sign    personal_sign    ${param}
    log    ${result}
    ${signature1} =    Get Substring    ${result}    2
    Set Global Variable    ${signature}    ${signature1}

isPulledPacket
    [Arguments]    ${tokenHolderPubKey}    ${message}
    ${param}    Create List    isPulledPacket    ${tokenHolderPubKey}    ${message}
    ${two}    Create List    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}    ${10}
    ${res}    post    contract_ccquery    isPulledPacket    ${two}
    log    ${res}
    [Return]    ${res}

getAllPacketInfo
    ${param}    Create List    getAllPacketInfo
    ${two}    Create List    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}    ${10}
    ${res}    post    contract_ccquery    getAllPacketInfo    ${two}
    log    ${res}
    ${addressMap}    To Json    ${res}
    log    ${addressMap}

createToken
    [Arguments]    ${address}    ${name}
    ${one}    Create List    createToken    BlackListTest    ${name}    1    1000
    ...    ${address}
    ${two}    Create List    ${address}    ${address}    0    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43
    ...    ${one}
    ${result}    post    contract_ccinvoketx    createToken    ${two}
    [Return]    ${result}

ccquery
    [Arguments]    ${name}
    ${one}    Create List    getTokenInfo    ${name}
    ${two}    Create List    PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43    ${one}    ${0}
    ${result}    post    contract_ccquery    getTokenInfo    ${two}
    ${addressMap}    To Json    ${result}
    ${assetId}    Get From Dictionary    ${addressMap}    AssetID
    [Return]    ${assetId}

createMultiTokenPacket
    [Arguments]    ${addr}    ${amount}    ${pubkey}    ${count}    ${min}    ${max}
    ...    ${expiredTime}    ${isConstant}    ${token1}    ${token2}
    ${param}    Create List    createPacket    ${pubkey}    ${count}    ${min}    ${max}
    ...    ${expiredTime}    remark    ${isConstant}
    ${two}    Create List    ${addr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${token1}    ${token2}    ${amount}
    ...    ${amount}    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}
    ${res}    post    contract_ccinvokeMutiToken    createPacket    ${two}
    log    ${res}    #    #    Create List    createPacket    ${pubkey}
    ...    # ${count}    ${min}    ${max}    # ${expiredTime}    remark    #
    ...    # Create List    ${addr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    PTN    ${amount}    1
    ...    # PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}
    [Return]    ${res}

multiTokenUpdated
    [Arguments]    ${addr}    ${token1}    ${amount1}    ${token2}    ${amount2}
    ${param}    Create List    updatePacket    ${tokenHolderPubKey}    0    0    0
    ...    ${EMPTY}    remark    true
    ${two}    Create List    ${addr}    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${token1}    ${token2}    ${amount1}
    ...    ${amount2}    1    PCGTta3M4t3yXu8uRgkKvaWd2d8DSDC6K99    ${param}
    ${res}    post    contract_ccinvokeMutiToken    updatePacket    ${two}
    log    ${res}
