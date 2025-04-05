{"author":"meet","post_dir":"thoughts","published":"published"}

<p>I was looking at a long list of logs, (debugging of course).</p>
<p>I had a list of transactions</p>
<ul>
<li>Two lists of transactions, one as ground truth and other as predicted.  Transactions is a list of object(dictionary or map), each object has fields like date, amount, description, etc.</li>
<li>Those objects in a list need not necessarily be in order, however I do want to compare them, how to do that?</li>
<li>I decided to sort them based on date(one field in that object).</li>
</ul>
<p>Then for each date, I group the transactions and this would narrow down the search space for comparison of one-one transaction, since now I can compare which of the ones are closely matching, the date will be a exact match the amount should also be, the description can be fuzzily matched.</p>
<p>However, for amount, I guess I was wrong, there could be a value like <code>10</code> and the other could have <code>10.01</code> and those python doesn't count equal atleast when compared as a string. I converted to float and compared rounded off numbers.</p>
<p>Now the problem kicked in and print statements flooded, dates everywhere.</p>
<p>Now I was using ghostty terminal, and there was definitely some scroll limit so I couldn't scroll all the way to start. So I opened up zellij.</p>
<p>I had a log with <em>MATCHED</em> text where I logged the date and amount where both transaction matched.</p>
<p>Now, I wanted a count of these matches. I could use search but it was not giving the count of occurrences, I can't keep counting with my mouth(that idea flew by though)</p>
<p>Now, that's where I accidentally hit <code>&lt;Ctrl&gt;S</code> and <code>E</code>
And I was in a editor, woah!</p>
<p>I was excited, I could finally copy and throw that in vim and get everything I want.
hehe
I tried +&quot;y but it didn't copy to the clipboard, it yanked yes, but not in the system clipboard. That frustrated me and took my hope down, but I googled it and also gpted it (is that a word, I think we can say llmed it).</p>
<p>And yes we can set the <code>export $EDITOR</code> to the editor executable path and it would open that thing in that editor.</p>
<p>I did and it did work.</p>
<p>That mode is called <a href="https://zellij.dev/news/edit-scrollback-compact/">scrollback-edit</a> mode. I should say a life saver mode, a log viewer and really cool.</p>
<p>I am probably too dumb and I know this exists in tmux, but I felt good and helped me solve my problem. So thank you whoever made that mode, its really helpful to debug with logs (debloging) Yes I am bad at naming things, but I like this more than vibe coding ;)</p>
