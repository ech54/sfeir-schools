pragma solidity ^0.4.22;

contract Bank {
 uint8 private clientCount; // define a counter for bank clients.
  address public owner; // define the owner of contract. type is address
  event WhoAmI(address who, uint256 balance);
  mapping (address => uint) private balances;  

  constructor() public payable {
      // set default value.
      emit WhoAmI(msg.sender, msg.sender.balance);
      clientCount = 0;
      owner = msg.sender; // msg allows to have details from Tx.
  }
    //deposit
    function deposit(address adr, uint value) public payable returns (uint) {
        balances[adr] += value;
        emit WhoAmI(adr, value);
        return balances[adr];
    }
    //withdraw
    function withdraw(address from, address to, uint withdrawAmount) public returns (uint remainingBal) {
            emit WhoAmI(from, balances[from]);
            emit WhoAmI(to, balances[to]);
            if (withdrawAmount <= balances[from]) {
                balances[from] -= withdrawAmount;
                to.transfer(withdrawAmount);
            }
            return balances[from];
        }
    // balance
    function balance(address adr) public constant returns (uint) {
        return balances[adr];
    }
}    


