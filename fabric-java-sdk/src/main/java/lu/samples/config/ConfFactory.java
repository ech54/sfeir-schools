package lu.samples.config;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * @TODO
 *
 * @author Emilien Charton
 * @date 14 nov. 2018
 */
public class ConfFactory {

	final Map<String, String> props = new HashMap<>();
	
	public static ConfFactory conf() {
		return new ConfFactory();
	}
	
	public ConfFactory prop(final String key, final String value) {
		this.props.put(key, value);
		return this;
	}
	
	public Conf instance() {
		return new Conf(props);
	}
	
	public class Conf {
		
		private final Map<String, String> props = new HashMap<>();
		
		private Conf(final Map<String, String> source) {
			this.props.putAll(source);
		}
		
		public Optional<String> val(final String k) {
			
			if (k==null || k.length()==0 
					|| !this.props.containsKey(k)) {
				return Optional.empty();
			}
			return Optional.of(this.props.get(k));
		}
		
	}
	
}
