package lu.samples.client;

import org.hyperledger.fabric.sdk.HFClient;
import org.junit.BeforeClass;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @TODO
 *
 * @author Emilien Charton
 * @date 14 nov. 2018
 */
public class ChannelTest {
	
	final Logger LOG = LoggerFactory.getLogger(ChannelTest.class);
	
	static HFClient client = null;
	
	@BeforeClass
	public static void mockFabricClient() {
		
		client = HFClient.createNewInstance();
		
	}
	
	
	@Test
	public void simpleConnect() {
		
		LOG.info("connect on channel: {} {}", "one","two");;
		
	}
	
}
