package lu.samples.scenario;

import java.io.File;
import java.net.MalformedURLException;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.Properties;
import java.util.stream.Stream;

import org.hyperledger.fabric.protos.peer.Query.ChaincodeInfo;
import org.hyperledger.fabric.sdk.BlockInfo;
import org.hyperledger.fabric.sdk.BlockchainInfo;
import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.HFClient;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.QueryByChaincodeRequest;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.EnrollmentRequest;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.protobuf.ByteString;

import lu.samples.config.Settings;
import lu.samples.model.UserDefault;

/**
 * @TODO
 *
 * @author Emilien Charton
 * @date 14 nov. 2018
 */
public class SimpleConnection {

	// Reference on the default logging system.
	final Logger LOG = LoggerFactory.getLogger(SimpleConnection.class);
	
	/**
	 * TODO
	 * 
	 * @param args
	 */
	public static void main(final String[] args) {
		new SimpleConnection().execute();
	}
	
	
	public void execute() {
		try {
			final UserDefault adminUser = fetchDomainAdmin(Settings.ORG1_USER_ADMIN, Settings.ORG1_CA);
			final HFClient client = instanciateFabricClient();
			client.setUserContext(adminUser);
			
			final Channel channel = fetchChannel(client);
			
			final QueryByChaincodeRequest request = client.newQueryProposalRequest();
			
			final ChaincodeID CCId = ChaincodeID.newBuilder().setName("sfeircc").setVersion("1.0").build();
		    request.setChaincodeID(CCId);
		    request.setFcn("read");
		    request.setArgs(new String[] { "0001" });
		    final Collection<ProposalResponse> res = channel.queryByChaincode(request);
		    
		    for (final ProposalResponse pres : res) {
				LOG.info("\n {}\n {}\n response-status: {} \n payload: {} \n endorsement: {} \n", 
						pres.getPeer(),
						pres.getChaincodeID(),
						pres.getChaincodeActionResponseStatus(), 
						new String(pres.getChaincodeActionResponsePayload()),
						pres.getProposalResponse().getEndorsement()
						);
		    }
			
		} 
		catch(final Exception exception) {
			LOG.error("", exception);
		}
	}
	
	protected void requestProposalOnChaincode(final Channel channel) {
		
	}
	
	protected Channel fetchChannel(final HFClient client) {
		Channel channel = null;
		try {
			final Peer peer = client.newPeer(Settings.ORG1_PEER0.val("name").get(), 
											 Settings.ORG1_PEER0.val("url").get());
//	        EventHub eventHub = client.newEventHub("eventhub01", "grpc://localhost:7053");
	        final Orderer orderer = client.newOrderer(Settings.ORDERER.val("name").get(), 
	        										  Settings.ORDERER.val("url").get());
	        channel = client.newChannel(Settings.CHANNEL.val("name").get());
	        channel.addPeer(peer);
	        
//	        channel.addEventHub(eventHub);
	        channel.addOrderer(orderer);
	        channel.initialize();
		} 
		catch(final Exception exception) {
			LOG.error("", exception);
		}
		
		return channel;
	}
	
	/**
	 * TODO
	 * @param userAdmin
	 * @param ca
	 * @return
	 */
	protected UserDefault fetchDomainAdmin(final Settings userAdmin, final Settings ca) {
		
		final UserDefault admin = new UserDefault();
		buildCaClient(ca.val("url").get(), clientProperties())
			.ifPresent(caCli -> {
				try {
					final String uName = userAdmin.val("name").get();
	                final EnrollmentRequest enrollmentRequestTLS = new EnrollmentRequest();
	                enrollmentRequestTLS.addHost("localhost");
	                enrollmentRequestTLS.setProfile("tls");
					final Enrollment enrollment = caCli.enroll("admin", userAdmin.val("pwd").get(), enrollmentRequestTLS);
					admin.setName(uName);
					admin.setDomain(userAdmin.val("domain").get());
					admin.setMspId(userAdmin.val("mspId").get());
					admin.setEnrollment(enrollment);
				} 
				catch(final Exception exception) {
					LOG.error("", exception);
				}
			});
		return admin;
	}
	
	/**
	 * Instantiate a default <code>HFClient</code> class. The instance is configured
	 *  with a default <code>CryptoSuite</code>.
	 * @return The new <code>HFClient</code> instance.
	 */
	protected HFClient instanciateFabricClient() {
		
		final HFClient client = HFClient.createNewInstance();
		instanciateCryptoSuite().ifPresent(crypto -> setupClientCrypto(client, crypto));
		return client;
	}
	
	// Build a Certificate Authoritative client baser on its URL.
	private Optional<HFCAClient> buildCaClient(final String url, final Properties props) {
		
		try {
			final HFCAClient caClient = HFCAClient.createNewInstance(url, props);
			instanciateCryptoSuite().ifPresent(crypto -> caClient.setCryptoSuite(crypto));
			return Optional.of(caClient);
		} 
		catch(final MalformedURLException exception) {
			LOG.error("", exception);
		}
		return Optional.empty();
	}
	
	// Instanciate the crypto suite <code>CryptoSuite</code>.
	private Optional<CryptoSuite> instanciateCryptoSuite() {
		try  {
			final CryptoSuite crypto = CryptoSuite.Factory.getCryptoSuite();
			return Optional.of(crypto);
		} 
		catch(final Exception exception) {
			LOG.error("Can't create cryptography configuration on fabric client.", exception);
		}
		
		return Optional.empty();
	}
	
	 // Setup the cryptography on the <code>HFClient</code> instance.
	private void setupClientCrypto(final HFClient target, final CryptoSuite crypto) {
		try {
			target.setCryptoSuite(crypto);
		} 
		catch(final Exception exception){
			LOG.error("Can't setup cryptography configuration on fabric client.", exception);
		}
	}
	
	private String defaultPath = "/home/echarton/data/home/wks/learning/tp/bc/sfeir-school/hyperledger/hyperledger_bc_intro-resolved/";
	private Properties clientProperties() {
        Properties ret = new Properties();
        final String domainName = "warehouse.sfeir.lu";
        final String componentType = "peer";
        File cert = Paths.get(defaultPath+"crypto-config/peerOrganizations/",domainName, componentType+"s", componentType+"0."+domainName, "/tls/server.crt").toFile();
        if (!cert.exists()) {
            throw new RuntimeException(String.format("Missing cert file for: %s. Could not find at location: %s", domainName,
                    cert.getAbsolutePath()));
        }
        File clientCert = Paths.get(defaultPath, "crypto-config/peerOrganizations/", domainName, "users/Admin@" + domainName, "tls/client.crt").toFile();
        File clientKey = Paths.get(defaultPath, "crypto-config/peerOrganizations/", domainName, "users/Admin@" + domainName, "tls/client.key").toFile();
        if (!clientCert.exists()) {
            throw new RuntimeException(String.format("Missing  client cert file for: %s. Could not find at location: %s", domainName, clientCert.getAbsolutePath()));
        }

        if (!clientKey.exists()) {
            throw new RuntimeException(String.format("Missing  client key file for: %s. Could not find at location: %s", domainName, clientKey.getAbsolutePath()));
        }
        ret.setProperty("clientCertFile", clientCert.getAbsolutePath());
        ret.setProperty("clientKeyFile", clientKey.getAbsolutePath());

//        ret.setProperty("pemFile", cert.getAbsolutePath());
//
        ret.setProperty("hostnameOverride", componentType+"0."+domainName);
        ret.setProperty("sslProvider", "openSSL");
        ret.setProperty("negotiationType", "TLS");
        return ret;
	}
}
