package config

func LogDefautConfig() string {
	return `
        <seelog type="sync">
        	<outputs formatid="main">
				<rollingfile type="size" filename="blingbling.log" maxsize="1073741824" maxrolls="10" />
        	</outputs>
            <formats>
                <format id="main" format="%Date %Time [%Level] %File:%Line %Msg%n"/>
            </formats>
        </seelog>
    `
}
