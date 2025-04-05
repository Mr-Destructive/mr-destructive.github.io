{"author":"meet","post_dir":"newsletter","published":"published"}

<h2>Week #36</h2>
<p>Was quite a work-holic week, did a lot of things this week. Debugging, prompting, scripting, flying through a random codebase and number crunching. I felt good, I am genuinely curious about Agents now and want to experiment a lot this weekend. Will have a video or two about it.</p>
<p>Finally it was a week, a little slow down after a 2 back to back weeks of AI launching and breaking the expectations. The image generation, gemma models, gemini 2.5 and mistral OCR, and what not have been released in the last stretch of first quarter of 2025. Phew! it was a heck of a quarter, that went in a jiffy.</p>
<h3>Quote of the week</h3>
<blockquote>
<p>“Its hard to beat talent when talent is driven by curiosity, but impossible to beat talent when it is driven by hard work“</p>
<p>— Meet Gor</p>
</blockquote>
<p>Believe me, talent and hard work is deadly combo. You don’t want to be competing with that kind of group of people, you might easily slip out, if you are just hard working or just talent won’t be enough, there needs to be a fire, a purpose, a curiosity burning that will lift even the one without feet and help fly without wings.</p>
<hr>
<h2>Created</h2>
<ul>
<li>
<p>Pydanitc-ai-compatible interface for Meta AI LLM (llama)</p>
<p>I created a pydantic class that will take the Meta AI LLM API and use it as a model instead of an API. I also need to make it tool-callable, since the api stucture is not designed for tool-calling, it needs to be done in a hacky way. But if that works, I can make the LLM available as an agent for free to anyone via the Meta AI API. That is a revolution.</p>
</li>
<li>
<p>API for converting a image URL into base 64 string representation in Golang and Javascript → https://imgbase.korogi5431.workers.dev/?url=</p>
<p>I created this to understand the process of encoding the image blob to base 64, first I get the Image URL, then fetch that image, sometimes, the Image might be encrypted or locked in due to authentication, so it will fail, however if the image is available publicly we can fetch it and load the bytes to encoded in base 64. This way, if the resource i.e. URI breaks or in future can’t reference it, the image link will be broken, can’t view the image. With the base 64 encoded value, the image is constructed in place, so we are not locked in with the network call. There are trade-offs, the string is quite big!</p>
</li>
</ul>
<h2>Read</h2>
<ul>
<li>
<p><a href="https://www.piglei.com/articles/en-programmer-reading-list-part-one/?ref=dailydev">A</a> <a href="https://www.piglei.com/articles/en-programmer-reading-list-part-one/?ref=dailydev">programmer’s</a> <a href="https://www.piglei.com/articles/en-programmer-reading-list-part-one/?ref=dailydev">Reading List: 100 Articles I Enjoyed Part 1</a> (1-50)Have not read any of those, but seems a great place to bookmark for never reading (bookmarks are bad design actually in 2025)</p>
</li>
<li>
<p><a href="https://open.substack.com/pub/zaidesanton/p/the-13-software-engineering-laws">The 13 software engineering laws</a>: This is a great read, full of controversial sometimes. Deadlines, answering wrong is the better than asking the wrong question, that was funny.</p>
</li>
</ul>
<h2>Watched</h2>
<ul>
<li><a href="https://youtu.be/3yrAK2hMWw8?si=rOnFCMbfbGn-wU_j"><em>I Ranked every LLM based on vibes</em></a></li>
</ul>
<p>{% embed https://youtu.be/3yrAK2hMWw8?si=rOnFCMbfbGn-wU_j %}</p>
<p>That is really cool, I didn’t knew gemini flash was that good, need to check that out really. I haven’t tried o3 and not sure would be quick enough to get results out. I have been using claude and gpt extensively for qick fixes and even brainstorming ideas. Pretty good tier list to be honest.</p>
<ul>
<li>
<p><a href="https://youtu.be/CDjjaTALI68?si=Dkse1dPU96MherDQ">Understanding MCP from scratch</a></p>
<p>This cleared a lot of things.</p>
</li>
<li>
<ul>
<li>
<p>RAG = Context + LLMs</p>
<ul>
<li>
<ul>
<li>
<p>Agents =Tools + LLMs</p>
<ul>
<li>MCP Server/Client = Context + Tools + LLM</li>
</ul>
</li>
</ul>
</li>
</ul>
</li>
</ul>
</li>
<li>
<p>{%embed https://youtu.be/CDjjaTALI68?si=Dkse1dPU96MherDQ %}</p>
</li>
<li>
<p><a href="https://youtu.be/Cor-t9xC1ck?si=8SuaiF8kHdI8hrQw">Raising an Agent - Episode 1</a></p>
<p>Fascinating. A sourcegraph AI Agent is around the corner. Wow this could be the first one to actually replicate the editing experience. We know LLMs can’t really drive the code, so let them do the chore work while we think. That is the approach Sourcegraph will be taking, they are not completely saying LLMs are bad, they are infact bullish on LLMs, claude can do almost anything provided the tools, so LLMs with tools and context is a big brain move.</p>
</li>
</ul>
<p>{% embed https://youtu.be/Cor-t9xC1ck?si=8SuaiF8kHdI8hrQw %}</p>
<h2>Learnt</h2>
<ul>
<li>
<p>How to convert a image URL to base 64 image (obviously loading the actual image bytes first)</p>
<p>So, we need to convert the bytes into base64 encoded letters (A-Z,a-z,0-9, and + and /). This is really cool to represent blobs of bits, this reduces the space of storing images, however its quite big for high quality images, as it will have a lot of information, like very granular pixel details.</p>
</li>
<li>
<p>Creating a model wrapper of Meta AI API: for pydantic compatible api</p>
</li>
<li>
<p>Creating a Cloud function on Appwrite Cloud in Golang</p>
</li>
<li>
<p>Creating a Cloud function worker on Cloudflare in Javascript</p>
</li>
</ul>
<h2>Tech News</h2>
<ul>
<li>
<p><a href="https://turso.tech/blog/turso-offline-sync-public-beta"><em>Turso</em> launches offline sync</a>: Sync your data to the primary database next time whenever you have the internet connectivity. This will allow the data to be saved on the local copy(replica) and get pushed to the primary database which will sync the primary and the replica. Really exited to try it out.</p>
</li>
<li>
<p><a href="https://www.picolm.io/">Pico LM</a>: Demystifying how Language models learn. This is really cool, I want to understand how it is trained and what improves or what deteriorates the accuracy.</p>
</li>
<li>
<p><a href="https://www.anthropic.com/news/introducing-claude-for-education">Anthropic launches program for Universities and Students</a> - Claude for Education</p>
<p>This is adoption, survival of the fittest, finding the ways in which people will reap benefits. Anthropic is solid step ahead of all the LLM providers in that race too, it is already solid at tool-calling, and now its spreading its legs.</p>
</li>
</ul>
<p>For more interesting articles, check out the <a href="https://buttondown.com/hacker-newsletter/archive/hacker-newsletter-740/">hacker-newsletter</a> for the week edition <a href="https://buttondown.com/hacker-newsletter/archive/hacker-newsletter-740/">#740</a>, for even more software development/coding articles, join daily.dev.</p>
<hr>
<p>A promising week ahead, lot of experiments to be done, a lot to be recorded, don’t have the mental energy but I will and need to push through to make myself better. I am really curious about AI Agents, I am emphasizing this again, because this is a field that needs attention. Because Attention is all you need!</p>
<p>See you next time :)</p>
<p>Thanks for reading <a href="https://techstructively.substack.com/"><em>Techstructive Weekly</em></a>! This post is public so feel free to share it.</p>
<p>Thanks for reading <a href="https://techstructively.substack.com/"><em>Techstructive Weekly</em></a>! Subscribe for free to receive new posts and support my work.</p>
