<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import CardSelect from '$lib/components/CardSelect.svelte';
	import Dropdown from '$lib/components/Dropdown.svelte';
	import Hand from '$lib/components/Hand.svelte';
	import Table from '$lib/components/Table.svelte';
	import WidowExchange from '$lib/components/WidowExchange.svelte';
	import { bidStore, gameStore, lobbyStore, userStore } from '$lib/stores';
	import type { BidInfo, Card, GameInfo } from '$lib/types';
	import { onDestroy, onMount } from 'svelte';
	
    $userStore = {
        name: "admin",
        admin: true,
        expireTime: 0,
    };
    $lobbyStore = {
        id: "test",
        host: "admin",
        players: ["admin", "dummy1", "dummy2", "dummy3"],
        started: false,
    };
    $bidStore = {
        id: "test",
        done: false,
        players: ["admin", "dummy1", "dummy2", "dummy3"],
        hand: [{color: 0, value: 0}, {color: 1, value: 1}, {color: 2, value: 2}, {color: 1, value: 12}],
        passed: [false, false, false, false],
        currentBidder: 0,
        bid: 0,
    };
    $gameStore = {
        id: "test",
        done: false,
        players: ["admin", "dummy1", "dummy2", "dummy3"],
        hand: [{color: 0, value: 0}, {color: 1, value: 1}, {color: 2, value: 2}, {color: 1, value: 12}],
		discardSize: [0, 0],
        table: [],
        currentPlayer: 0,
        trump: 0,
        bid: 100,
        bidWinner: 0,
    };

    async function getWidow() {
        return [{color: 1, value: 2}, {color: 1, value: 8}, {color: 2, value: 3}, {color: 4, value: 13}];
    }

    async function startRound(trumpSelection: number, toWidow: Card[], fromWidow: Card[]) {
        if ($gameStore === undefined) {
            return;
        }
        $gameStore.trump = trumpSelection;
    }

	const yourIndex = $gameStore?.players.indexOf($userStore?.name ?? '') ?? -1;
	const tookBid = $gameStore?.bidWinner === yourIndex;

	$: trumpSelected = $gameStore?.trump ?? 0 != 0;
	$: yourHand = $gameStore?.hand ?? [];

	let trumpSelection: number = 0;

	$: currentPlayer = $gameStore?.currentPlayer ?? -1;

	let trumpColor = '';
	$: if (trumpSelected) {
		trumpColor = getTrumpColor($gameStore);
	}
	function getTrumpColor(gameInfo: GameInfo | undefined) {
		if (gameInfo === undefined) {
			return '';
		}
		const trump = gameInfo.trump;
		if (trump === 1) {
			return 'Red';
		}
		if (trump === 2) {
			return 'Yellow';
		}
		if (trump === 3) {
			return 'Green';
		}
		if (trump === 4) {
			return 'Black';
		}
		return '';
	}

	let toWidow: Card[] = [];
	let fromWidow: Card[]  = [];
	let startGameStatus = "";
	async function submitCreateGame() {
		// if (toWidow.length !== fromWidow.length) {
		// 	startGameStatus = "You must take the same number of cards out of your hand as you take out of the widow.";
		// 	return;
		// }
		// if (trumpSelection === 0) {
		// 	startGameStatus = "You need to choose the trump color.";
		// 	return;
		// }
		// const [ok, status] = await startRound(trumpSelection, toWidow, fromWidow);
		// if (!ok) {
		// 	console.log("When trying to start round, got status = " + status);
		// }
        await startRound(trumpSelection, toWidow, fromWidow);
	}

	let selectedCard: Card;
	async function submitSelectCard() {
		if ($gameStore === undefined) {return;}
		if (selectedCard === undefined) {
			return;
		}
		$gameStore.table.push(selectedCard);
		$gameStore.hand = $gameStore.hand.filter((v) => v !== selectedCard);
		// const [ok, status] = await playCard(selectedCard);
		// if (!ok) {
		// 	console.log("When trying to play card, got status = " + status);
		// }
	}
</script>
