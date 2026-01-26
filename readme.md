# MAGI Token

ERC-20 compliant token contract for the Magi Network.

## Overview

MAGI Token is a fungible token implementation with standard ERC-20 functionality plus additional features for pausability, ownership management, and safe allowance operations.

## Token Configuration

Token properties are configured at initialization via the `init` payload:

| Property   | Type   | Description                          |
|------------|--------|--------------------------------------|
| name       | string | Token name (e.g., "Magi Token")      |
| symbol     | string | Token symbol (e.g., "MAGI")          |
| decimals   | uint8  | Decimal places (e.g., 3)             |
| maxSupply  | uint64 | Maximum mintable supply              |

Example init payload:
```json
{"name": "Magi Token", "symbol": "MAGI", "decimals": 3, "maxSupply": 1000000000}
```

## Features

- **ERC-20 Compliant**: Standard token interface (`transfer`, `transferFrom`, `approve`, `allowance`, `balanceOf`, `totalSupply`)
- **Mintable**: Owner can mint new tokens up to max supply
- **Burnable**: Any holder can burn their own tokens
- **Pausable**: Owner can pause/unpause transfers
- **Ownership Transfer**: Owner can transfer contract ownership
- **Safe Allowance**: `increaseAllowance`/`decreaseAllowance` to prevent race conditions

## Functions

### Actions (State-Changing)

| Function            | Payload                                      | Access    |
|---------------------|----------------------------------------------|-----------|
| `init`              | `{"name": string, "symbol": string, "decimals": uint8, "maxSupply": uint64}` | ContractOwner |
| `mint`              | `{"amount": uint64}`                         | Owner     |
| `burn`              | `{"amount": uint64}`                         | Any       |
| `transfer`          | `{"to": string, "amount": uint64}`           | Any       |
| `transferFrom`      | `{"from": string, "to": string, "amount": uint64}` | Any |
| `approve`           | `{"spender": string, "amount": uint64}`      | Any       |
| `increaseAllowance` | `{"spender": string, "amount": uint64}`      | Any       |
| `decreaseAllowance` | `{"spender": string, "amount": uint64}`      | Any       |
| `pause`             | -                                            | Owner     |
| `unpause`           | -                                            | Owner     |
| `changeOwner`       | `{"newOwner": string}`                       | Owner     |

### Queries (Read-Only)

| Function      | Payload                                | Response                     |
|---------------|----------------------------------------|------------------------------|
| `balanceOf`   | `{"account": string}`                  | `{"balance": uint64}`        |
| `totalSupply` | -                                      | `{"totalSupply": uint64}`    |
| `allowance`   | `{"owner": string, "spender": string}` | `{"allowance": uint64}`      |
| `getOwner`    | -                                      | `{"owner": string}`          |
| `getInfo`     | -                                      | `{"name", "symbol", "decimals", "maxSupply"}` |
| `isPaused`    | -                                      | `{"paused": bool}`           |

## Events

All events include `type`, `attributes`, and `tx` (transaction ID).

| Event Type         | Attributes                                    |
|--------------------|-----------------------------------------------|
| `init_magi_token`  | `owner`, `name`, `symbol`, `decimals`, `maxSupply` |
| `transfer`         | `from`, `to`, `amount`                        |
| `approval`         | `owner`, `spender`, `amount`                  |
| `ownerChange`      | `previousOwner`, `newOwner`                   |
| `paused`           | `by`                                          |
| `unpaused`         | `by`                                          |

### ERC-20 Event Compliance

- **Mint**: Emits `transfer` with `from: ""`
- **Burn**: Emits `transfer` with `to: ""`
- **transferFrom**: Emits both `transfer` and `approval` (for updated allowance)

## Allowance Pattern (DEX Integration)

The allowance mechanism allows a third party (like a DEX) to transfer tokens on behalf of a user. This is essential for decentralized exchanges and other DeFi protocols.

### How It Works

```
1. User approves DEX to spend their tokens
2. DEX can transfer tokens from user to any recipient (up to approved amount)
3. Each transfer reduces the allowance
```

### DEX Integration Steps

**Step 1: User Approves DEX**

The user must first approve the DEX contract to spend their tokens:

```json
// User calls approve
{
  "action": "approve",
  "payload": {"spender": "hive:dex_contract", "amount": 5000000}
}
```

**Step 2: DEX Executes Transfer**

When a trade occurs, the DEX transfers tokens from the user:

```json
// DEX contract calls transferFrom
{
  "action": "transferFrom",
  "payload": {"from": "hive:user", "to": "hive:buyer", "amount": 1000000}
}
```

**Step 3: Allowance Auto-Decrements**

After the transfer, the allowance is automatically reduced:
- Original allowance: 5,000,000
- Transfer amount: 1,000,000
- Remaining allowance: 4,000,000

### Safe Allowance Management

To prevent race conditions (front-running attacks), use `increaseAllowance` and `decreaseAllowance` instead of setting absolute values with `approve`:

```json
// Increase existing allowance by 1000
{"action": "increaseAllowance", "payload": {"spender": "hive:dex", "amount": 1000}}

// Decrease existing allowance by 500
{"action": "decreaseAllowance", "payload": {"spender": "hive:dex", "amount": 500}}

// Revoke all allowance (set to 0)
{"action": "approve", "payload": {"spender": "hive:dex", "amount": 0}}
```

### Example: Complete DEX Swap Flow

```
1. Alice has 10,000 MAGI tokens
2. Alice approves DEX for 5,000 tokens
3. Bob wants to buy 2,000 MAGI from Alice
4. DEX calls transferFrom(alice, bob, 2000)
5. Result:
   - Alice: 8,000 MAGI
   - Bob: 2,000 MAGI
   - Alice's DEX allowance: 3,000
6. DEX can still transfer up to 3,000 more from Alice
```

### Important Notes for DEX Developers

- Always check `allowance` before attempting `transferFrom`
- `transferFrom` will fail if allowance < amount
- `transferFrom` will fail if owner's balance < amount
- Both `transfer` and `approval` events are emitted on `transferFrom`
- Allowance management (`approve`, `increaseAllowance`, `decreaseAllowance`) works even when contract is paused
- Actual transfers (`transfer`, `transferFrom`) are blocked when paused

## Build

```bash
tinygo build -gc=custom -scheduler=none -panic=trap -no-debug -target=wasm-unknown -o test/artifacts/main.wasm ./contract
```

## Test

```bash
go test ./test/...
```

## Project Structure

```
magi_token/
├── contract/
│   ├── main.go            # Entry point and state helpers
│   ├── token.go           # Exported WASM functions
│   ├── internal.go        # Internal helper functions
│   ├── types.go           # Type definitions
│   ├── types_tinyjson.go  # JSON serialization (tinyjson)
│   └── events.go          # Event emission
├── sdk/                   # VSC SDK bindings
├── test/
│   ├── basic_test.go      # Core token tests (init, mint, transfer, burn)
│   ├── allowance_test.go  # Allowance tests (approve, transferFrom)
│   ├── pausable_test.go   # Pause/unpause tests
│   ├── edge_cases_test.go # Edge cases & negative tests
│   ├── lifecycle_test.go  # Comprehensive lifecycle test
│   ├── helpers_test.go    # Test utilities
│   └── artifacts/         # Compiled WASM
└── readme.md
```

## RC Consumption

| Function           | Avg RC |
|--------------------|--------|
| Queries            | 100    |
| unpause            | 110    |
| pause              | 122    |
| burn               | 170    |
| changeOwner        | 165    |
| decreaseAllowance  | 191    |
| approve            | 205    |
| increaseAllowance  | 196    |
| transfer           | 252    |
| mint               | 326    |
| transferFrom       | 443    |
| init               | 948    |
