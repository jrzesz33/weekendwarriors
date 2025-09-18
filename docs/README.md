# Weekend Warrior Game Tracking

## Description

For the bogey golfer and high handicapper, this application keeps track of scores and calculates live score tracking for side bet games during a golf round. 

## Input

The app should allow an anoymous user to create a new game and input a name and handicap for each golfer. The app should also guide users using information like a bogey golfer for men is a 20 handicap, and 24 for women if they are unsure what to put in. The app should also provide them a link for them to share with other players to follow along and see the scores. Once they start the game an interface should be allowed for them to input each players score and the number of putts, for each hole starting with the first hole. The app should also show the active score for the side bets happening.

## Side Bets

The app should handle two different side bets.
- Best Nine
    The app should calculate the score compared to par with handicap based on the best nine holes of the round. If a player has nine holes shooting +3 that were the worst holes they shot, and the remaining holes equaled +18, +18 would be that players score. If handicap is enabled the app should calculate that into the score against the other players. 
- Putt Putt Poker
    Each player starts with three playing cards. The app should get putting strokes for each score. If a player gets a one putt, they get an additional card. If they zero-putt, whole in, they get two cards. If they three putt or worst they should be notified to add a dollar to the bet. Handicap does not apply with this bet. 

## End of Round

At the end of the round and after all of the players scores for all 18 holes are submitted, the app should randomly deal the playing cards and calculate the players hands with the best poker hand. The app should also display the winner of the Best Nine Game.

## Non-Functional Requirements

- The app should use the Diamond Run golf course for where this is being played and here is the holes and pars. (Hole and Par)
    1. 4	
    2. 3	
    3. 5	
    4. 4	
    5. 3	
    6. 4
    7. 4
    8. 3
    9. 5
    10. 4
    11. 3
    12. 5
    13.	3
    14.	4
    15.	4
    16.	4
    17.	5
    18.	4